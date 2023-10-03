package orderrepositorybun

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (r *OrderRepositoryBun) UpdateDeliveryOrder(ctx context.Context, order *orderentity.Order, delivery *orderentity.DeliveryOrder) error {
	return nil
}

func (r *OrderRepositoryBun) DeleteDeliveryOrder(ctx context.Context, id string) error {
	return nil
}
