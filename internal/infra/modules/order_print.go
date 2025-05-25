package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"

	orderprintusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_print"
)

// NewOrderPrintModule registers endpoints for printing individual orders.
func NewOrderPrintModule(db *bun.DB, chi *server.ServerChi) (*orderprintusecases.Service, *handler.Handler) {
	printSvc := orderprintusecases.NewService()

	handler := handlerimpl.NewHandlerOrderPrint(printSvc)
	chi.AddHandler(handler)
	return printSvc, handler
}
