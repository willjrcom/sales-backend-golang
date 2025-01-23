package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	tablerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/table"
	tableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table"
)

func NewTableModule(db *bun.DB, chi *server.ServerChi) (model.TableRepository, *tableusecases.Service, *handler.Handler) {
	repository := tablerepositorybun.NewTableRepositoryBun(db)
	service := tableusecases.NewService(repository)
	handler := handlerimpl.NewHandlerTable(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
