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

func NewOrderModule(db *bun.DB, chi *server.ServerChi, stockRepo model.StockRepository, stockMovementRepo model.StockMovementRepository) (model.OrderRepository, *orderusecases.OrderService, *handler.Handler) {
	repository := orderrepositorybun.NewOrderRepositoryBun(db)
	service := orderusecases.NewOrderService(repository)

	// Adicionar dependências de estoque
	service.AddDependencies(
		repository,
		nil, // rs - será injetado depois
		nil, // rp - será injetado depois
		nil, // rpr - será injetado depois
		nil, // rdo - será injetado depois
		stockRepo,
		stockMovementRepo,
		nil, // sgi - será injetado depois
		nil, // sop - será injetado depois
		nil, // sq - será injetado depois
		nil, // sd - será injetado depois
		nil, // sp - será injetado depois
		nil, // st - será injetado depois
		nil, // sc - será injetado depois
	)

	handler := handlerimpl.NewHandlerOrder(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
