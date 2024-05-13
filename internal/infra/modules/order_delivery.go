package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	orderdeliveryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_delivery"
)

func NewOrderDeliveryModule(db *bun.DB, chi server.ServerChi) (*orderrepositorybun.OrderDeliveryRepositoryBun, orderdeliveryusecases.IService, *handler.Handler) {
	repository := orderrepositorybun.NewOrderDeliveryRepositoryBun(db)
	service := orderdeliveryusecases.NewService(repository)
	handler := handlerimpl.NewHandlerOrderDelivery(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
