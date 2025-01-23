package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	orderprocessrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order_process"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

func NewOrderProcessModule(db *bun.DB, chi *server.ServerChi) (model.OrderProcessRepository, *orderusecases.OrderProcessService, *handler.Handler) {
	repository := orderprocessrepositorybun.NewOrderProcessRepositoryBun(db)
	service := orderusecases.NewOrderProcessService(repository)
	handler := handlerimpl.NewHandlerOrderProcess(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
