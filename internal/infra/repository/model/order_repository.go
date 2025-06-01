package model

import (
	"context"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	PendingOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	GetAllOrdersWithDelivery(ctx context.Context, page, perPage int) ([]Order, error)
	AddPaymentOrder(ctx context.Context, payment *PaymentOrder) error
}
