package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"

	orintmanagerusecases "github.com/willjrcom/sales-backend-go/internal/usecases/print_manager"
)

// NewOrderPrintModule registers endpoints for printing individual orders.
func NewOrderPrintModule(db *bun.DB, chi *server.ServerChi) (*orintmanagerusecases.Service, *handler.Handler) {
	printSvc := orintmanagerusecases.NewService()

	handler := handlerimpl.NewHandlerPrintManager(printSvc)
	chi.AddHandler(handler)
	return printSvc, handler
}
