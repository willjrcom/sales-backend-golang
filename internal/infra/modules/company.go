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
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

func NewCompanyModule(db *bun.DB, chi *server.ServerChi) (model.CompanyRepository, *companyusecases.Service, *handler.Handler) {
	repository := companyrepositorybun.NewCompanyRepositoryBun(db)
	mpClient := mercadopagoservice.NewClient()
	service := companyusecases.NewService(repository, mpClient)
	service.StartSubscriptionWatcher(context.Background(), 24*time.Hour)
	handler := handlerimpl.NewHandlerCompany(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
