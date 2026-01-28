package modules

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/scheduler"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
	billingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

func NewCompanyModule(db *bun.DB, chi *server.ServerChi, costRepo model.CompanyUsageCostRepository) (model.CompanyRepository, *companyusecases.Service, *handler.Handler) {
	repository := companyrepositorybun.NewCompanyRepositoryBun(db)
	companyPaymentRepo := companyrepositorybun.NewCompanyPaymentRepositoryBun(db)
	mpClient := mercadopagoservice.NewClient()
	service := companyusecases.NewService(repository, companyPaymentRepo)
	service.StartSubscriptionWatcher(context.Background(), 24*time.Hour)

	checkoutUC := billingusecases.NewCheckoutUseCase(costRepo, repository, companyPaymentRepo, mpClient)
	costService := companyusecases.NewUsageCostService(costRepo, repository)

	// Start Monthly Billing Scheduler
	billingScheduler := scheduler.NewMonthlyBillingScheduler(checkoutUC, repository, companyPaymentRepo)
	billingScheduler.Start(context.Background())

	handler := handlerimpl.NewHandlerCompany(service, checkoutUC, costService, billingScheduler)
	chi.AddHandler(handler)
	return repository, service, handler
}
