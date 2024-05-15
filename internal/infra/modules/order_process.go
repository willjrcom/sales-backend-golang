package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	orderprocessrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order_process"
	orderprocessusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_process"
)

func NewOrderProcessModule(db *bun.DB, chi *server.ServerChi) (*orderprocessrepositorybun.ProcessRepositoryBun, *orderprocessusecases.Service, *handler.Handler) {
	repository := orderprocessrepositorybun.NewOrderProcessRepositoryBun(db)
	service := orderprocessusecases.NewService(repository)
	handler := handlerimpl.NewHandlerOrderProcess(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
