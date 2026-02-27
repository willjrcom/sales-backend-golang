package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	billing "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type DailyScheduler struct {
	db                      *bun.DB
	companyRepo             model.CompanyRepository
	orderRepo               model.OrderRepository
	companyPaymentRepo      model.CompanyPaymentRepository
	companySubscriptionRepo model.CompanySubscriptionRepository
	checkoutUseCase         *billing.CheckoutUseCase
	companyUseCase          *companyusecases.Service
	orderUseCase            *orderusecases.OrderService
}

func NewDailyScheduler(db *bun.DB, companyRepo model.CompanyRepository, orderRepo model.OrderRepository, companyPaymentRepo model.CompanyPaymentRepository, companySubscriptionRepo model.CompanySubscriptionRepository, checkoutUseCase *billing.CheckoutUseCase, companyUseCase *companyusecases.Service, orderUseCase *orderusecases.OrderService) *DailyScheduler {
	return &DailyScheduler{
		db:                      db,
		companyRepo:             companyRepo,
		orderRepo:               orderRepo,
		companyPaymentRepo:      companyPaymentRepo,
		companySubscriptionRepo: companySubscriptionRepo,
		checkoutUseCase:         checkoutUseCase,
		companyUseCase:          companyUseCase,
		orderUseCase:            orderUseCase,
	}
}

func (s *DailyScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case t := <-ticker.C:
				// Run billing checks at 5 AM
				if t.Hour() == 5 {
					log.Println("Running Daily Batch...")
					s.ProcessCostsToPay(ctx)
					s.UpdateCompanyPlans(ctx)
					s.CheckOverdueAccounts(ctx)
					s.CheckExpiredOptionalPayments(ctx)
					s.CleanStagingOrders(ctx)
					log.Println("Daily Batch Completed.")
				}
			}
		}
	}()
}

func (s *DailyScheduler) CleanStagingOrders(ctx context.Context) {
	log.Println("Scheduler: Cleaning up staging orders...")
	var schemas []string
	if err := s.db.NewRaw("SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname LIKE 'company_%'").Scan(ctx, &schemas); err != nil {
		log.Printf("Scheduler: Error fetching schemas: %v", err)
		return
	}

	for _, schema := range schemas {
		ctxSchema := context.WithValue(ctx, model.Schema("schema"), schema)

		orders, err := s.orderRepo.GetOrdersByStatus(ctxSchema, orderentity.OrderStatusStaging)
		if err != nil {
			log.Printf("Scheduler: Error fetching staging orders in schema %s: %v", schema, err)
			continue
		}

		for _, order := range orders {
			dtoID := &entitydto.IDRequest{ID: order.ID}
			if err := s.orderUseCase.CancelOrder(ctxSchema, dtoID, true); err != nil {
				log.Printf("Scheduler: Error cancelling staging order %s in schema %s: %v", order.ID, schema, err)
			}
		}

		log.Printf("Scheduler: Cleaned %d staging orders in schema %s", len(orders), schema)
	}
}

func (s *DailyScheduler) UpdateCompanyPlans(ctx context.Context) error {
	return s.companySubscriptionRepo.UpdateCompanyPlans(ctx)
}

func (s *DailyScheduler) CheckOverdueAccounts(ctx context.Context) {
	// 1. Block companies with overdue payments (> 5 days)
	cutoffDate := time.Now().UTC().AddDate(0, 0, -5)
	overduePayments, err := s.companyPaymentRepo.ListOverduePayments(ctx, cutoffDate)
	if err == nil {
		for _, payment := range overduePayments {
			_ = s.companyRepo.UpdateBlockStatus(ctx, payment.CompanyID, true)
		}
	}

	// 2. Unblock companies that have settled their mandatory payments
	companies, err := s.companyRepo.ListBlockCompaniesForBilling(ctx)
	if err == nil {
		for _, company := range companies {
			if company.IsBlocked {
				pending, err := s.companyPaymentRepo.ListPendingMandatoryPayments(ctx, company.ID)
				if err == nil && len(pending) == 0 {
					_ = s.companyRepo.UpdateBlockStatus(ctx, company.ID, false)
				}
			}
		}
	}
}

func (s *DailyScheduler) ProcessCostsToPay(ctx context.Context) {
	// Calculate target date (10 days from now)
	targetDate := time.Now().UTC().AddDate(0, 0, 10)
	targetDay := targetDate.Day()

	daysProcessing := []int{targetDay}

	// Check if targetDate is the last day of the month
	if targetDate.AddDate(0, 0, 1).Month() != targetDate.Month() {
		// If last day, include all subsequent days (e.g. if 28th Feb, include 29, 30, 31)
		for d := targetDay + 1; d <= 31; d++ {
			daysProcessing = append(daysProcessing, d)
		}
	}

	for _, day := range daysProcessing {
		companies, err := s.companyRepo.ListCompaniesByPaymentDueDay(ctx, day)
		if err != nil {
			// log error
			continue
		}

		for _, company := range companies {
			_ = s.checkoutUseCase.GenerateMonthlyCostPayment(ctx, company.ID)
		}
	}
}

func (s *DailyScheduler) CheckExpiredOptionalPayments(ctx context.Context) {
	payments, err := s.companyPaymentRepo.ListExpiredOptionalPayments(ctx)
	if err != nil {
		// log error
		return
	}

	for _, payment := range payments {
		// Reuse CheckoutUseCase.CancelPayment logic (unlinks costs, updates status)
		// Assuming CancelPayment handles idempotency or allowed status checks
		_ = s.companyUseCase.CancelPayment(ctx, payment.ID)
	}
}
