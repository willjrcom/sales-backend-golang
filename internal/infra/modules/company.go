package modules

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/scheduler"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/focusnfe"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
	billingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

func NewCompanyModule(db *bun.DB, chi *server.ServerChi, costRepo model.CompanyUsageCostRepository) (model.CompanyRepository, *companyusecases.Service, *handler.Handler) {
	repository := companyrepositorybun.NewCompanyRepositoryBun(db)
	companyPaymentRepo := companyrepositorybun.NewCompanyPaymentRepositoryBun(db)
	mpClient := mercadopagoservice.NewClient()
	focusClient := focusnfe.NewClient()
	service := companyusecases.NewService(repository, companyPaymentRepo, focusClient)
	// service.StartSubscriptionWatcher removed in favor of DailyScheduler

	checkoutUC := billingusecases.NewCheckoutUseCase(costRepo, repository, companyPaymentRepo, mpClient)
	costService := companyusecases.NewUsageCostService(costRepo, repository)

	// Start Daily Scheduler
	dailyScheduler := scheduler.NewDailyScheduler(repository, companyPaymentRepo, checkoutUC, service)
	dailyScheduler.Start(context.Background())

	handler := handlerimpl.NewHandlerCompany(service, checkoutUC, costService, dailyScheduler)
	chi.AddHandler(handler)
	return repository, service, handler
}
