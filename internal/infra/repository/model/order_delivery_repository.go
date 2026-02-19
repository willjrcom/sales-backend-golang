package model

import "context"

type OrderDeliveryRepository interface {
	CreateOrderDelivery(ctx context.Context, delivery *OrderDelivery) error
	UpdateOrderDelivery(ctx context.Context, delivery *OrderDelivery) error
	DeleteOrderDelivery(ctx context.Context, id string) error
	GetDeliveryById(ctx context.Context, id string) (*OrderDelivery, error)
	GetDeliveriesByIds(ctx context.Context, ids []string) ([]OrderDelivery, error)
	GetOrderIDFromOrderDeliveriesByClientId(ctx context.Context, clientID string) ([]OrderDelivery, error)
	GetAllDeliveries(ctx context.Context) ([]OrderDelivery, error)
}
