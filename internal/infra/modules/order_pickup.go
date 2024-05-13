package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	orderpickupusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_pickup"
)

func NewOrderPickupModule(db *bun.DB, chi server.ServerChi) (*orderrepositorybun.OrderPickupRepositoryBun, orderpickupusecases.IService, *handler.Handler) {
	repository := orderrepositorybun.NewOrderPickupRepositoryBun(db)
	service := orderpickupusecases.NewService(repository)
	handler := handlerimpl.NewHandlerOrderPickup(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
