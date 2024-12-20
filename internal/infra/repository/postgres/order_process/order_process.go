package orderprocessrepositorybun

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type ProcessRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderProcessRepositoryBun(db *bun.DB) *ProcessRepositoryBun {
	return &ProcessRepositoryBun{db: db}
}

func (r *ProcessRepositoryBun) CreateProcess(ctx context.Context, s *orderprocessentity.OrderProcess) error {
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
		tx.Rollback()
		return err
	}

	if len(s.Products) != 0 {
		for _, p := range s.Products {
			processToProduct := &orderprocessentity.OrderProcessToProductToGroupItem{
				ProcessID:   s.ID,
				ProductID:   p.ID,
				GroupItemID: s.GroupItemID,
			}

			if _, err := tx.NewInsert().Model(processToProduct).Exec(ctx); err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					continue
				}

				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProcessRepositoryBun) UpdateProcess(ctx context.Context, s *orderprocessentity.OrderProcess) error {
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

	if _, err := r.db.NewDelete().Model(&orderprocessentity.OrderProcess{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcessRepositoryBun) GetProcessById(ctx context.Context, id string) (*orderprocessentity.OrderProcess, error) {
	process := &orderprocessentity.OrderProcess{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(process).Where("process.id = ?", id).
		Relation("GroupItem.Items.AdditionalItems").
		Relation("GroupItem.ComplementItem").
		Relation("GroupItem.Category").Scan(ctx); err != nil {
		return nil, err
	}

	return process, nil
}

func (r *ProcessRepositoryBun) GetAllProcesses(ctx context.Context) ([]orderprocessentity.OrderProcess, error) {
	processes := []orderprocessentity.OrderProcess{}

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

func (r *ProcessRepositoryBun) GetProcessesByProcessRuleID(ctx context.Context, id string) ([]orderprocessentity.OrderProcess, error) {
	processes := []orderprocessentity.OrderProcess{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processes).Where("process.process_rule_id = ? and process.status != ?", id, orderprocessentity.ProcessStatusFinished).
		Relation("GroupItem.Items.AdditionalItems").
		Relation("GroupItem.ComplementItem").
		Relation("GroupItem.Category").
		Scan(ctx); err != nil {
		return nil, err
	}

	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByProductID(ctx context.Context, id string) ([]orderprocessentity.OrderProcess, error) {
	processes := []orderprocessentity.OrderProcess{}
	processesToProduct := []orderprocessentity.OrderProcessToProductToGroupItem{}

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

func (r *ProcessRepositoryBun) GetProcessesByGroupItemID(ctx context.Context, id string) ([]orderprocessentity.OrderProcess, error) {
	processes := []orderprocessentity.OrderProcess{}
	processesToGroupItem := []orderprocessentity.OrderProcessToProductToGroupItem{}

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
