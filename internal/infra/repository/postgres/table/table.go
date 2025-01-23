package tablerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type TableRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewTableRepositoryBun(db *bun.DB) model.TableRepository {
	return &TableRepositoryBun{db: db}
}

func (r *TableRepositoryBun) CreateTable(ctx context.Context, s *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableRepositoryBun) UpdateTable(ctx context.Context, s *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableRepositoryBun) DeleteTable(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.Table{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TableRepositoryBun) GetTableById(ctx context.Context, id string) (*model.Table, error) {
	table := &model.Table{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(table).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return table, nil
}

func (r *TableRepositoryBun) GetAllTables(ctx context.Context) ([]model.Table, error) {
	tables := make([]model.Table, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&tables).Scan(ctx); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *TableRepositoryBun) GetUnusedTables(ctx context.Context) ([]model.Table, error) {
	tables := make([]model.Table, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&tables).Where("id NOT IN (?)", r.db.NewSelect().Model((*model.PlaceToTables)(nil)).
		Column("table_id")).Scan(ctx); err != nil {
		return nil, err
	}

	return tables, nil
}
