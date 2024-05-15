package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	deliverydriverrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/delivery_driver"
	deliverydriverusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_driver"
)

func NewDeliveryDriverModule(db *bun.DB, chi *server.ServerChi) (*deliverydriverrepositorybun.DeliveryDriverRepositoryBun, *deliverydriverusecases.Service, *handler.Handler) {
	repository := deliverydriverrepositorybun.NewDeliveryDriverRepositoryBun(db)
	service := deliverydriverusecases.NewService(repository)
	handler := handlerimpl.NewHandlerDeliveryDriver(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
