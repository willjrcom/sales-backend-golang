package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	employeerepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/employee"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
)

func NewEmployeeModule(db *bun.DB, chi *server.ServerChi) (*employeerepositorybun.EmployeeRepositoryBun, *employeeusecases.Service, *handler.Handler) {
	repository := employeerepositorybun.NewEmployeeRepositoryBun(db)
	service := employeeusecases.NewService(repository)
	handler := handlerimpl.NewHandlerEmployee(service)
	chi.AddHandler(handler)
	return repository, service, handler
}
