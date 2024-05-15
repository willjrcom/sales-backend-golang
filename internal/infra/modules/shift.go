package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	shiftrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/shift"
	shiftusecases "github.com/willjrcom/sales-backend-go/internal/usecases/shift"
)

func NewShiftModule(db *bun.DB, chi *server.ServerChi) (*shiftrepositorybun.ShiftRepositoryBun, *shiftusecases.Service, *handler.Handler) {
	repository := shiftrepositorybun.NewShiftRepositoryBun(db)
	service := shiftusecases.NewService(repository)
	handler := handlerimpl.NewHandlerShift(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
