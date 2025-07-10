package model

import "context"

type StockRepository interface {
	CreateStock(ctx context.Context, s *Stock) error
	UpdateStock(ctx context.Context, s *Stock) error
	GetStockByID(ctx context.Context, id string) (*Stock, error)
	GetStockByProductID(ctx context.Context, productID string) (*Stock, error)
	GetAllStocks(ctx context.Context) ([]Stock, error)
	GetActiveStocks(ctx context.Context) ([]Stock, error)
	GetLowStockProducts(ctx context.Context) ([]Stock, error)
	GetOutOfStockProducts(ctx context.Context) ([]Stock, error)
}

type StockMovementRepository interface {
	CreateMovement(ctx context.Context, m *StockMovement) error
	GetMovementsByStockID(ctx context.Context, stockID string) ([]StockMovement, error)
	GetMovementsByProductID(ctx context.Context, productID string) ([]StockMovement, error)
	GetMovementsByOrderID(ctx context.Context, orderID string) ([]StockMovement, error)
	GetAllMovements(ctx context.Context) ([]StockMovement, error)
	GetMovementsByDateRange(ctx context.Context, start, end string) ([]StockMovement, error)
}

type StockAlertRepository interface {
	CreateAlert(ctx context.Context, a *StockAlert) error
	UpdateAlert(ctx context.Context, a *StockAlert) error
	GetAlertByID(ctx context.Context, alertID string) (*StockAlert, error)
	GetAllAlerts(ctx context.Context) ([]StockAlert, error)
	GetAlertsByStockID(ctx context.Context, stockID string) ([]StockAlert, error)
	GetActiveAlerts(ctx context.Context) ([]StockAlert, error)
	GetResolvedAlerts(ctx context.Context) ([]StockAlert, error)
	ResolveAlert(ctx context.Context, alertID string, resolvedBy string) error
	DeleteAlert(ctx context.Context, alertID string) error
}
