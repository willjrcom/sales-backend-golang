package processrepositorybun

import (
	"context"
	"database/sql"
	"strings"
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

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if len(s.Products) != 0 {
		for _, p := range s.Products {
			processToProduct := &processentity.ProcessToProductToGroupItem{
				ProcessID:   s.ID,
				ProductID:   p.ID,
				GroupItemID: s.GroupItemID,
			}

			if _, err := tx.NewInsert().Model(processToProduct).Exec(ctx); err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					continue
				}

				if errRollBack := tx.Rollback(); errRollBack != nil {
					return errRollBack
				}

				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
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

	if err := r.db.NewSelect().Model(&processes).Scan(ctx); err != nil {
		return nil, err
	}

	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByProductID(ctx context.Context, id string) ([]processentity.Process, error) {
	processes := []processentity.Process{}
	processesToProduct := []processentity.ProcessToProductToGroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processesToProduct).Where("product_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	processIDs := []string{}
	for _, p := range processesToProduct {
		processIDs = append(processIDs, p.ProcessID.String())
	}

	if err := r.db.NewSelect().Model(&processes).Where("id in (?)", bun.In(processIDs)).Scan(ctx); err != nil {
		return nil, err
	}

	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByGroupItemID(ctx context.Context, id string) ([]processentity.Process, error) {
	processes := []processentity.Process{}
	processesToGroupItem := []processentity.ProcessToProductToGroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processesToGroupItem).Where("group_item_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	processIDs := []string{}
	for _, p := range processesToGroupItem {
		processIDs = append(processIDs, p.ProcessID.String())
	}

	if err := r.db.NewSelect().Model(&processes).Where("id in (?)", bun.In(processIDs)).Scan(ctx); err != nil {
		return nil, err
	}

	return processes, nil
}
