package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	userrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/user"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
)

func NewUserModule(db *bun.DB, chi server.ServerChi) (*userrepositorybun.UserRepositoryBun, *userusecases.Service, *handler.Handler) {
	repository := userrepositorybun.NewUserRepositoryBun(db)
	service := userusecases.NewService(repository)
	handler := handlerimpl.NewHandlerUser(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
