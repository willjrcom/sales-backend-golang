package model

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	PendingOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAllOrders(ctx context.Context, shiftID string, withStatus []orderentity.StatusOrder, withCategory bool) ([]Order, error)
	GetAllOrdersWithDelivery(ctx context.Context, shiftID string, page, perPage int) ([]Order, error)
	GetAllOrdersWithPickup(ctx context.Context, shiftID string, page, perPage int) ([]Order, error)
	AddPaymentOrder(ctx context.Context, payment *PaymentOrder) error
}
