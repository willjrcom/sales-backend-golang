package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/company"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

func NewCompanyModule(db *bun.DB, chi *server.ServerChi) (model.CompanyRepository, *companyusecases.Service, *handler.Handler) {
	repository := companyrepositorybun.NewCompanyRepositoryBun(db)
	service := companyusecases.NewService(repository)
	handler := handlerimpl.NewHandlerCompany(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
