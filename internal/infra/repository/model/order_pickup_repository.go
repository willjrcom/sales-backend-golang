package model

import "context"

type OrderPickupRepository interface {
	CreateOrderPickup(ctx context.Context, pickup *OrderPickup) error
	UpdateOrderPickup(ctx context.Context, pickup *OrderPickup) error
	DeleteOrderPickup(ctx context.Context, id string) error
	GetPickupById(ctx context.Context, id string) (*OrderPickup, error)
	GetOrderIDFromOrderPickupsByContact(ctx context.Context, contact string) ([]OrderPickup, error)
	GetAllPickups(ctx context.Context) ([]OrderPickup, error)
}
