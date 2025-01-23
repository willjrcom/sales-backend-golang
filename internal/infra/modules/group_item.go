package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	groupitemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/group_item"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

func NewGroupItemModule(db *bun.DB, chi *server.ServerChi) (model.GroupItemRepository, *orderusecases.GroupItemService, *handler.Handler) {
	repository := groupitemrepositorybun.NewGroupItemRepositoryBun(db)
	service := orderusecases.NewGroupItemService(repository)
	handler := handlerimpl.NewHandlerGroupItem(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
