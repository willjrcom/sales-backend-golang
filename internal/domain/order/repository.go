package orderentity

import "context"

type Repository interface {
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	UpdateDeliveryOrder(ctx context.Context, order *Order, delivery *DeliveryOrder) error
	UpdateTableOrder(ctx context.Context, order *Order, table *TableOrder) error
	DeleteOrder(ctx context.Context, id string) error
	DeleteDeliveryOrder(ctx context.Context, id string) error
	DeleteTableOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrderBy(ctx context.Context, o *Order) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
}
