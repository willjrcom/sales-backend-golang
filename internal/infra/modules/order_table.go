package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	orderrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/order"
	ordertableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_table"
)

func NewOrderTableModule(db *bun.DB, chi server.ServerChi) (*orderrepositorybun.OrderTableRepositoryBun, *ordertableusecases.Service, *handler.Handler) {
	repository := orderrepositorybun.NewOrderTableRepositoryBun(db)
	service := ordertableusecases.NewService(repository)
	handler := handlerimpl.NewHandlerOrderTable(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
