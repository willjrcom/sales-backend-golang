package orderrepositorybun

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (r *OrderRepositoryBun) UpdateTableOrder(ctx context.Context, order *orderentity.Order, table *orderentity.TableOrder) error {
	return nil
}

func (r *OrderRepositoryBun) DeleteTableOrder(ctx context.Context, id string) error {
	return nil
}
