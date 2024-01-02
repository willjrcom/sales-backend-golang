package tablerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type TableRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewTableRepositoryBun(db *bun.DB) *TableRepositoryBun {
	return &TableRepositoryBun{db: db}
}

func (r *TableRepositoryBun) RegisterTable(ctx context.Context, s *tableentity.Table) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepositoryBun) UpdateTable(ctx context.Context, s *tableentity.Table) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepositoryBun) DeleteTable(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&tableentity.Table{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *TableRepositoryBun) GetTableById(ctx context.Context, id string) (*tableentity.Table, error) {
	table := &tableentity.Table{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(table).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return table, nil
}
