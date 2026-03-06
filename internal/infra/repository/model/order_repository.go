package model

import (
	"context"

	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	PendingOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	UpdateOrderWithRelations(ctx context.Context, order *Order) error
	UpdateOrderWithRelationsWithTx(ctx context.Context, tx *bun.Tx, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	DeleteOrdersByStatus(ctx context.Context, status orderentity.StatusOrder) error
	ExistsOrderById(ctx context.Context, id string) (bool, error)
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrderByIdWithTx(ctx context.Context, tx *bun.Tx, id string) (*Order, error)
	GetOnlyOrderById(ctx context.Context, id string) (*Order, error)
	GetAllOpenedOrders(ctx context.Context) ([]Order, error)
	GetStaleStagingOrders(ctx context.Context, minutes int) ([]Order, error)
	GetAllOrders(ctx context.Context, shiftID string, withStatus []orderentity.StatusOrder) ([]Order, error)
	GetAllOrdersWithReadyDelivery(ctx context.Context, page, perPage int) ([]Order, error)
	GetAllOrdersWithShippedDelivery(ctx context.Context, page, perPage int) ([]Order, error)
	GetAllOrdersWithFinishedDelivery(ctx context.Context, shiftID string, page, perPage int) ([]Order, error)
	GetAllOrdersWithPickup(ctx context.Context, status orderentity.StatusOrderPickup, page, perPage int) ([]Order, error)
	GetOrdersByStatus(ctx context.Context, status orderentity.StatusOrder) ([]Order, error)
	AddPaymentOrder(ctx context.Context, payment *PaymentOrder) error
}
