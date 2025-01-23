package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	itemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/item"
	itemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/item"
)

func NewItemModule(db *bun.DB, chi *server.ServerChi) (model.ItemRepository, *itemusecases.Service, *handler.Handler) {
	repository := itemrepositorybun.NewItemRepositoryBun(db)
	service := itemusecases.NewService(repository)
	handler := handlerimpl.NewHandlerItem(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
