package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderTableRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderTableRepositoryBun(db *bun.DB) *OrderTableRepositoryBun {
	return &OrderTableRepositoryBun{db: db}
}

func (r *OrderTableRepositoryBun) CreateOrderTable(ctx context.Context, table *model.OrderTable) error {
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

func (r *OrderTableRepositoryBun) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {
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

func (r *OrderTableRepositoryBun) DeleteOrderTable(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.OrderTable{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderTableRepositoryBun) GetOrderTableById(ctx context.Context, id string) (table *model.OrderTable, err error) {
	table = &model.OrderTable{}
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

func (r *OrderTableRepositoryBun) GetPendingOrderTablesByTableId(ctx context.Context, id string) (tables []model.OrderTable, err error) {
	tables = []model.OrderTable{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&tables).Where("table_id = ? AND status = ?", id, orderentity.OrderTableStatusPending).Scan(ctx); err != nil {
		return nil, err
	}

	return tables, err
}

func (r *OrderTableRepositoryBun) GetAllOrderTables(ctx context.Context) (tables []model.OrderTable, err error) {
	tables = make([]model.OrderTable, 0)

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
