package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
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
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(table).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableOrderRepositoryBun) UpdateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(table).WherePK().Where("id = ?", table.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableOrderRepositoryBun) DeleteTableOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&orderentity.TableOrder{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableOrderRepositoryBun) GetTableOrderById(ctx context.Context, id string) (table *orderentity.TableOrder, err error) {
	table = &orderentity.TableOrder{}
	table.ID = uuid.MustParse(id)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err = r.db.NewSelect().Model(table).WherePK().Scan(ctx); err != nil {
		return nil, err
	}

	return table, err
}

func (r *TableOrderRepositoryBun) GetPendingTableOrdersByTableId(ctx context.Context, id string) (tables []orderentity.TableOrder, err error) {
	tables = []orderentity.TableOrder{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&tables).Where("table_id = ? AND status = ?", id, orderentity.TableOrderStatusPending).Scan(ctx); err != nil {
		return nil, err
	}

	return tables, err
}

func (r *TableOrderRepositoryBun) GetAllTableOrders(ctx context.Context) (tables []orderentity.TableOrder, err error) {
	tables = make([]orderentity.TableOrder, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err = r.db.NewSelect().Model(&tables).Scan(ctx); err != nil {
		return nil, err
	}

	return tables, err
}
