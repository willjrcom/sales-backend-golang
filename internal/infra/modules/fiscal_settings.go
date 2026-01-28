package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	fiscalsettingsrepository "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/fiscal_settings"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/focusnfe"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	fiscalsettingsusecases "github.com/willjrcom/sales-backend-go/internal/usecases/fiscal_settings"
)

func NewFiscalSettingsModule(db *bun.DB, chi *server.ServerChi, companyRepo model.CompanyRepository, companyService *companyusecases.Service) *handler.Handler {
	fiscalSettingsRepo := fiscalsettingsrepository.NewFiscalSettingsRepositoryBun(db)
	focusClient := focusnfe.NewClient()
	service := fiscalsettingsusecases.NewService(fiscalSettingsRepo, companyRepo, focusClient)
	handler := handlerimpl.NewFiscalSettingsHandler(service)

	chi.AddHandler(handler)

	return handler
}
