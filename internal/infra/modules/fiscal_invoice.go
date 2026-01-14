package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	fiscalinvoicerepository "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/fiscal_invoice"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/transmitenota"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	fiscalinvoiceusecases "github.com/willjrcom/sales-backend-go/internal/usecases/fiscal_invoice"
)

// NewFiscalInvoiceModule initializes fiscal invoice and usage cost modules
func NewFiscalInvoiceModule(
	db *bun.DB,
	chi *server.ServerChi,
	companyRepo model.CompanyRepository,
	orderRepo model.OrderRepository,
	companyService *companyusecases.Service,
) (model.FiscalInvoiceRepository, *fiscalinvoiceusecases.Service, *companyusecases.UsageCostService) {
	// Repositories
	fiscalInvoiceRepo := fiscalinvoicerepository.NewFiscalInvoiceRepository(db)
	usageCostRepo := companyrepositorybun.NewCompanyUsageCostRepository(db)

	// Services
	transmitenotaClient := transmitenota.NewClient()
	usageCostService := companyusecases.NewUsageCostService(usageCostRepo, companyRepo)
	fiscalInvoiceService := fiscalinvoiceusecases.NewService(
		fiscalInvoiceRepo,
		companyRepo,
		orderRepo,
		usageCostService,
		transmitenotaClient,
	)

	// Handlers
	usageCostHandler := handlerimpl.NewHandlerCompanyUsageCost(usageCostService, companyService)
	fiscalInvoiceHandler := handlerimpl.NewHandlerFiscalInvoice(fiscalInvoiceService)

	chi.AddHandler(usageCostHandler)
	chi.AddHandler(fiscalInvoiceHandler)

	return fiscalInvoiceRepo, fiscalInvoiceService, usageCostService
}
