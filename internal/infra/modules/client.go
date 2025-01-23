package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	clientrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/client"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
)

func NewClientModule(db *bun.DB, chi *server.ServerChi) (model.ClientRepository, *clientusecases.Service, *handler.Handler) {
	repository := clientrepositorybun.NewClientRepositoryBun(db)
	service := clientusecases.NewService(repository)
	handler := handlerimpl.NewHandlerClient(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
