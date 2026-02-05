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

func NewCompanyModule(db *bun.DB, chi *server.ServerChi, costRepo model.CompanyUsageCostRepository) (model.CompanyRepository, *companyusecases.Service, *billingusecases.CheckoutUseCase, *handler.Handler) {
	companyRepository := companyrepositorybun.NewCompanyRepositoryBun(db)
	companySubscriptionRepo := companyrepositorybun.NewCompanySubscriptionRepositoryBun(db)
	companyPaymentRepo := companyrepositorybun.NewCompanyPaymentRepositoryBun(db)
	mpClient := mercadopagoservice.NewClient()
	focusClient := focusnfe.NewClient()
	service := companyusecases.NewService(companyRepository, companyPaymentRepo, focusClient)
	// service.StartSubscriptionWatcher removed in favor of DailyScheduler

	checkoutUC := billingusecases.NewCheckoutUseCase(costRepo, companyRepository, companyPaymentRepo, companySubscriptionRepo, mpClient)
	costService := companyusecases.NewUsageCostService(costRepo, companyRepository)

	// Start Daily Scheduler
	dailyScheduler := scheduler.NewDailyScheduler(companyRepository, companyPaymentRepo, companySubscriptionRepo, checkoutUC, service)
	dailyScheduler.Start(context.Background())

	handler := handlerimpl.NewHandlerCompany(service, checkoutUC, costService, dailyScheduler)
	chi.AddHandler(handler)
	return companyRepository, service, checkoutUC, handler
}
