package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	groupitemrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/group_item"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
)

func NewGroupItemModule(db *bun.DB, chi server.ServerChi) (*groupitemrepositorybun.GroupItemRepositoryBun, *groupitemusecases.Service, *handler.Handler) {
	repository := groupitemrepositorybun.NewGroupItemRepositoryBun(db)
	service := groupitemusecases.NewService(repository)
	handler := handlerimpl.NewHandlerGroupItem(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
