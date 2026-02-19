package model

import "context"

type OrderTableRepository interface {
	CreateOrderTable(ctx context.Context, table *OrderTable) error
	UpdateOrderTable(ctx context.Context, table *OrderTable) error
	DeleteOrderTable(ctx context.Context, id string) error
	GetOrderTableById(ctx context.Context, id string) (*OrderTable, error)
	GetPendingOrderTablesByTableId(ctx context.Context, id string) ([]OrderTable, error)
	GetOrderTablesByTableId(ctx context.Context, id string, contact string) ([]OrderTable, error)
	GetAllOrderTables(ctx context.Context) ([]OrderTable, error)
}
