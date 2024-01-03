package orderentity

import "context"

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	PendingOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
}

type DeliveryOrderRepository interface {
	CreateDeliveryOrder(ctx context.Context, delivery *DeliveryOrder) error
	UpdateDeliveryOrder(ctx context.Context, delivery *DeliveryOrder) error
	DeleteDeliveryOrder(ctx context.Context, id string) error
	GetDeliveryById(ctx context.Context, id string) (*DeliveryOrder, error)
	GetAllDeliveries(ctx context.Context) ([]DeliveryOrder, error)
}

type TableOrderRepository interface {
	CreateTableOrder(ctx context.Context, table *TableOrder) error
	UpdateTableOrder(ctx context.Context, table *TableOrder) error
	DeleteTableOrder(ctx context.Context, id string) error
	GetTableOrderById(ctx context.Context, id string) (*TableOrder, error)
	GetAllTableOrders(ctx context.Context) ([]TableOrder, error)
}
