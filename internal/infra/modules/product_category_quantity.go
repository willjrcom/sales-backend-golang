package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	quantityrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/quantity"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
)

func NewProductCategoryQuantityModule(db *bun.DB, chi *server.ServerChi) (model.QuantityRepository, *quantityusecases.Service, *handler.Handler) {
	repository := quantityrepositorybun.NewQuantityRepositoryBun(db)
	service, err := quantityusecases.InitializeService()
	if err != nil {
		panic(err)
	}

	handler := handlerimpl.NewHandlerQuantity(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
