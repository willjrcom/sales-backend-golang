package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	contactrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/contact"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact"
)

func NewContactModule(db *bun.DB, chi server.ServerChi) (*contactrepositorybun.ContactRepositoryBun, *contactusecases.Service, *handler.Handler) {
	repository := contactrepositorybun.NewContactRepositoryBun(db)
	service := contactusecases.NewService(repository)
	handler := handlerimpl.NewHandlerContactPerson(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
