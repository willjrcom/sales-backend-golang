package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

func NewOrderTableModule(db *bun.DB, chi *server.ServerChi) (model.OrderTableRepository, *orderusecases.OrderTableService, *handler.Handler) {
	repository := orderrepositorybun.NewOrderTableRepositoryBun(db)
	service := orderusecases.NewOrderTableService(repository)
	handler := handlerimpl.NewHandlerOrderTable(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
