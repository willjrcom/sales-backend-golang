package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	billing "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

type DailyScheduler struct {
	companyRepo             model.CompanyRepository
	companyPaymentRepo      model.CompanyPaymentRepository
	companySubscriptionRepo model.CompanySubscriptionRepository
	checkoutUseCase         *billing.CheckoutUseCase
	companyUseCase          *companyusecases.Service
}

func NewDailyScheduler(companyRepo model.CompanyRepository, companyPaymentRepo model.CompanyPaymentRepository, companySubscriptionRepo model.CompanySubscriptionRepository, checkoutUseCase *billing.CheckoutUseCase, companyUseCase *companyusecases.Service) *DailyScheduler {
	return &DailyScheduler{
		companyRepo:             companyRepo,
		companyPaymentRepo:      companyPaymentRepo,
		companySubscriptionRepo: companySubscriptionRepo,
		checkoutUseCase:         checkoutUseCase,
		companyUseCase:          companyUseCase,
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
				// Run billing checks at 8 AM
				if t.Hour() == 5 {
					log.Println("Running Daily Billing Batch...")
					s.ProcessCostsToPay(ctx)
					s.UpdateCompanyPlans(ctx)
					s.CheckOverdueAccounts(ctx)
					s.CheckExpiredOptionalPayments(ctx)
					log.Println("Daily Billing Batch Completed.")
				}
			}
		}
	}()
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
