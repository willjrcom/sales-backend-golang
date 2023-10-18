package orderrepositorybun

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (r *OrderRepositoryBun) UpdateDeliveryOrder(ctx context.Context, delivery *orderentity.DeliveryOrder) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(delivery).WherePK().Where("id = ?", delivery.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) DeleteDeliveryOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&orderentity.DeliveryOrder{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}
