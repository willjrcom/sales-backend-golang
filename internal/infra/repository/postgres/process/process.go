package processrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

type ProcessRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProcessRepositoryBun(db *bun.DB) *ProcessRepositoryBun {
	return &ProcessRepositoryBun{db: db}
}

func (r *ProcessRepositoryBun) RegisterProcess(ctx context.Context, s *processentity.Process) error {
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

func (r *ProcessRepositoryBun) UpdateProcess(ctx context.Context, s *processentity.Process) error {
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

func (r *ProcessRepositoryBun) DeleteProcess(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&processentity.Process{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcessRepositoryBun) GetProcessById(ctx context.Context, id string) (*processentity.Process, error) {
	process := &processentity.Process{}

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

func (r *ProcessRepositoryBun) GetAllProcesses(ctx context.Context) ([]processentity.Process, error) {
	processes := []processentity.Process{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(processes).Scan(ctx); err != nil {
		return nil, err
	}

	return processes, nil
}
