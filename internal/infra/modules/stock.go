package modules

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/server"
	handlerimpl "github.com/willjrcom/sales-backend-go/internal/infra/handler"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	stockrepositorybun "github.com/willjrcom/sales-backend-go/internal/infra/repository/postgres/stock"
	stockusecases "github.com/willjrcom/sales-backend-go/internal/usecases/stock"
)

func NewStockModule(db *bun.DB, chi *server.ServerChi, productRepo model.ProductRepository) (model.StockRepository, *stockusecases.Service, interface{}) {
	// Repositories
	stockRepo := stockrepositorybun.NewStockRepositoryBun(db)
	stockMovementRepo := stockrepositorybun.NewStockMovementRepositoryBun(db)
	stockAlertRepo := stockrepositorybun.NewStockAlertRepositoryBun(db)

	// Use cases
	stockService := stockusecases.NewStockService(stockRepo, stockMovementRepo, stockAlertRepo, productRepo)

	// Handlers
	stockHandler := handlerimpl.NewHandlerStock(stockService)

	// Add handler to server
	chi.AddHandler(stockHandler)

	return stockRepo, stockService, stockHandler
}
