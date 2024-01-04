package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
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

func (r *TableOrderRepositoryBun) GetTableOrderById(ctx context.Context, id string) (table *orderentity.TableOrder, err error) {
	table = &orderentity.TableOrder{}
	table.ID = uuid.MustParse(id)

	r.mu.Lock()
	err = r.db.NewSelect().Model(table).WherePK().Relation("Waiter").Scan(ctx)
	r.mu.Unlock()

	return table, err
}

func (r *TableOrderRepositoryBun) GetAllTableOrders(ctx context.Context) (tables []orderentity.TableOrder, err error) {
	tables = make([]orderentity.TableOrder, 0)

	r.mu.Lock()
	err = r.db.NewSelect().Model(&tables).Relation("Waiter").Scan(ctx)
	r.mu.Unlock()

	return tables, err
}
