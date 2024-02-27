package processrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProcessCategoryRepositoryBun(db *bun.DB) *ProcessCategoryRepositoryBun {
	return &ProcessCategoryRepositoryBun{db: db}
}

func (r *ProcessCategoryRepositoryBun) RegisterProcess(ctx context.Context, s *productentity.Process) error {
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

func (r *ProcessCategoryRepositoryBun) UpdateProcess(ctx context.Context, s *productentity.Process) error {
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

func (r *ProcessCategoryRepositoryBun) DeleteProcess(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&productentity.Process{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcessCategoryRepositoryBun) GetProcessById(ctx context.Context, id string) (*productentity.Process, error) {
	process := &productentity.Process{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(process).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return process, nil
}
