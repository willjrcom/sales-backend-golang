package orderentity

import "context"

type Repository interface {
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrderBy(ctx context.Context, o *Order) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
}
