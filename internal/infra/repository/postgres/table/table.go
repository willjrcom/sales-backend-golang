package tablerepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type TableRepositoryBun struct {
	db *bun.DB
}

func NewTableRepositoryBun(db *bun.DB) model.TableRepository {
	return &TableRepositoryBun{db: db}
}

func (r *TableRepositoryBun) CreateTable(ctx context.Context, s *model.Table) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryBun) UpdateTable(ctx context.Context, s *model.Table) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryBun) DeleteTable(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false
	isActive := false
	if _, err := tx.NewUpdate().
		Model(&model.Table{}).
		Set("is_active = ?", isActive).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryBun) GetTableById(ctx context.Context, id string) (*model.Table, error) {
	table := &model.Table{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(table).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return table, nil
}

func (r *TableRepositoryBun) GetAllTables(ctx context.Context, page, perPage int, isActive bool) ([]model.Table, int, error) {
	tables := make([]model.Table, 0)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	query := tx.NewSelect().Model(&tables).Where("is_active = ?", isActive).Limit(perPage).Offset(page * perPage)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return tables, count, nil
}

func (r *TableRepositoryBun) GetUnusedTables(ctx context.Context, isActive ...bool) ([]model.Table, error) {
	tables := make([]model.Table, 0)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Default to active records (true)
	activeFilter := true
	if len(isActive) > 0 {
		activeFilter = isActive[0]
	}

	if err := tx.NewSelect().Model(&tables).Where("id NOT IN (?) AND is_active = ?", tx.NewSelect().Model((*model.PlaceToTables)(nil)).
		Column("table_id"), activeFilter).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return tables, nil
}
