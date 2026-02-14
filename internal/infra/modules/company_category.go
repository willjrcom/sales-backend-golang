package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companycategoryrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company_category"
	companycategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company_category"
)

func NewCompanyCategoryModule(db *bun.DB, chi *server.ServerChi) (model.CompanyCategoryRepository, *companycategoryusecases.Service, *handler.Handler) {
	companyCategoryRepository := companycategoryrepositorybun.NewCompanyCategoryRepositoryBun(db)
	service := companycategoryusecases.NewService(companyCategoryRepository)
	handler := handlerimpl.NewHandlerCompanyCategory(service)
	chi.AddHandler(handler)
	return companyCategoryRepository, service, handler
}
