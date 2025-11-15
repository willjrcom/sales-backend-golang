package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	stockrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/stock"
	stockusecases "github.com/willjrcom/sales-backend-go/internal/usecases/stock"
)

func NewStockModule(db *bun.DB, chi *server.ServerChi) (model.StockRepository, model.StockMovementRepository, model.StockAlertRepository, *stockusecases.Service, *handler.Handler) {
	// Repositories
	stockRepo := stockrepositorybun.NewStockRepositoryBun(db)
	stockMovementRepo := stockrepositorybun.NewStockMovementRepositoryBun(db)
	stockAlertRepo := stockrepositorybun.NewStockAlertRepositoryBun(db)

	// Use cases
	stockService := stockusecases.NewStockService(stockRepo, stockMovementRepo, stockAlertRepo)

	// Handlers
	stockHandler := handlerimpl.NewHandlerStock(stockService)

	// Add handler to server
	chi.AddHandler(stockHandler)

	return stockRepo, stockMovementRepo, stockAlertRepo, stockService, stockHandler
}
