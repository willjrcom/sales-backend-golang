package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	orderqueuerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order_queue"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

func NewOrderqueueModule(db *bun.DB, chi *server.ServerChi) (*orderqueuerepositorybun.QueueRepositoryBun, *orderqueueusecases.Service, *handler.Handler) {
	repository := orderqueuerepositorybun.NewOrderQueueRepositoryBun(db)
	service := orderqueueusecases.NewService(repository)
	handler := handlerimpl.NewHandlerOrderQueue(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
