package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type TableOrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewTableOrderRepositoryBun(db *bun.DB) *TableOrderRepositoryBun {
	return &TableOrderRepositoryBun{db: db}
}

func (r *TableOrderRepositoryBun) CreateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(table).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *TableOrderRepositoryBun) UpdateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(table).WherePK().Where("id = ?", table.ID).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *TableOrderRepositoryBun) DeleteTableOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&orderentity.TableOrder{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	return err
}
