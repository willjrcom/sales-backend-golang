package orderprocessrepositorybun

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProcessRepositoryBun struct {
	db *bun.DB
}

func NewOrderProcessRepositoryBun(db *bun.DB) model.OrderProcessRepository {
	return &ProcessRepositoryBun{db: db}
}

func (r *ProcessRepositoryBun) CreateProcess(ctx context.Context, s *model.OrderProcess) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if len(s.Products) != 0 {
		for _, p := range s.Products {
			processToProduct := &model.OrderProcessToProductToGroupItem{
				ProcessID:   s.ID,
				ProductID:   p.ID,
				GroupItemID: s.GroupItemID,
			}

			if _, err := tx.NewInsert().Model(processToProduct).Exec(ctx); err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					continue
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

func (r *ProcessRepositoryBun) UpdateProcess(ctx context.Context, s *model.OrderProcess) error {

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

func (r *ProcessRepositoryBun) DeleteProcess(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.OrderProcess{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProcessRepositoryBun) GetProcessById(ctx context.Context, id string) (*model.OrderProcess, error) {
	process := &model.OrderProcess{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(process).Where("process.id = ?", id).
		Relation("GroupItem.Items.AdditionalItems").
		Relation("GroupItem.ComplementItem").
		Relation("GroupItem.Category").
		Relation("Products").
		Relation("Queue").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return process, nil
}

func (r *ProcessRepositoryBun) GetAllProcesses(ctx context.Context) ([]model.OrderProcess, error) {
	processes := []model.OrderProcess{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processes).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByProcessRuleID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	processes := []model.OrderProcess{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	validStatus := []orderprocessentity.StatusProcess{
		orderprocessentity.ProcessStatusFinished,
		orderprocessentity.ProcessStatusCanceled,
	}

	if err := tx.NewSelect().Model(&processes).
		Where("\"process\".process_rule_id = ? and \"process\".status NOT IN (?)",
			id, bun.In(validStatus)).
		Relation("GroupItem").
		Relation("GroupItem.Items", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("is_additional = ?", false)
		}).
		Relation("GroupItem.Items.AdditionalItems").
		Relation("GroupItem.ComplementItem").
		Relation("GroupItem.Category").
		Relation("Products").
		Relation("Queue").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByProductID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	processes := []model.OrderProcess{}
	processesToProduct := []model.OrderProcessToProductToGroupItem{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processesToProduct).Where("product_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	processIDs := []string{}
	for _, p := range processesToProduct {
		processIDs = append(processIDs, p.ProcessID.String())
	}

	if err := tx.NewSelect().Model(&processes).Where("id in (?)", bun.In(processIDs)).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processes, nil
}

func (r *ProcessRepositoryBun) GetProcessesByGroupItemID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	processes := []model.OrderProcess{}
	processesToGroupItem := []model.OrderProcessToProductToGroupItem{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processesToGroupItem).Where("group_item_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	processIDs := []string{}
	for _, p := range processesToGroupItem {
		processIDs = append(processIDs, p.ProcessID.String())
	}

	if err := tx.NewSelect().Model(&processes).Where("id in (?)", bun.In(processIDs)).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processes, nil
}
