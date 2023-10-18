package orderrepositorybun

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (r *OrderRepositoryBun) CreateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(table).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *OrderRepositoryBun) UpdateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(table).WherePK().Where("id = ?", table.ID).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *OrderRepositoryBun) DeleteTableOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&orderentity.TableOrder{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	return err
}
