package orderrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderTableRepositoryBun struct {
	db *bun.DB
}

func NewOrderTableRepositoryBun(db *bun.DB) model.OrderTableRepository {
	return &OrderTableRepositoryBun{db: db}
}

func (r *OrderTableRepositoryBun) CreateOrderTable(ctx context.Context, table *model.OrderTable) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(table).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderTableRepositoryBun) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(table).WherePK().Where("id = ?", table.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderTableRepositoryBun) DeleteOrderTable(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.OrderTable{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderTableRepositoryBun) GetOrderTableById(ctx context.Context, id string) (table *model.OrderTable, err error) {
	table = &model.OrderTable{}
	table.ID = uuid.MustParse(id)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err = tx.NewSelect().Model(table).WherePK().Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return table, err
}

func (r *OrderTableRepositoryBun) GetPendingOrderTablesByTableId(ctx context.Context, id string) (tables []model.OrderTable, err error) {
	tables = []model.OrderTable{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&tables).Where("table_id = ? AND status = ?", id, orderentity.OrderTableStatusPending).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return tables, err
}

func (r *OrderTableRepositoryBun) GetAllOrderTables(ctx context.Context) (tables []model.OrderTable, err error) {
	tables = make([]model.OrderTable, 0)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err = tx.NewSelect().Model(&tables).Where("status != 'Closed' AND status != 'Cancelled'").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return tables, err
}
