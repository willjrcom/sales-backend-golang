package processrulerepositorybun

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProcessRuleRepositoryBun struct {
	db *bun.DB
}

func NewProcessRuleRepositoryBun(db *bun.DB) model.ProcessRuleRepository {
	return &ProcessRuleRepositoryBun{db: db}
}

func (r *ProcessRuleRepositoryBun) CreateProcessRule(ctx context.Context, s *model.ProcessRule) error {

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

func (r *ProcessRuleRepositoryBun) UpdateProcessRule(ctx context.Context, s *model.ProcessRule) error {

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

func (r *ProcessRuleRepositoryBun) DeleteProcessRule(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	isActive := false
	if _, err := tx.NewUpdate().Model(&model.ProcessRule{}).Set("is_active = ?", isActive).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProcessRuleRepositoryBun) GetProcessRuleById(ctx context.Context, id string) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(processRule).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(processRule).Where("category_id = ? and order = 1", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRules := []model.ProcessRule{}
	mapProcesses := map[uuid.UUID]uuid.UUID{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processRules).Where("\"order\" = 1").Scan(ctx); err != nil {
		return nil, err
	}

	for _, processRule := range processRules {
		mapProcesses[processRule.CategoryID] = processRule.ID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return mapProcesses, nil
}

func (r *ProcessRuleRepositoryBun) GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(processRule).Where("category_id = ? and \"order\" = ?", id, order).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRulesSubquery := []model.ProcessRule{}
	processRules := []model.ProcessRule{}
	mapProcesses := map[uuid.UUID]uuid.UUID{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Subconsulta para obter o máximo order para cada category_id
	err = tx.NewSelect().Model(&processRulesSubquery).ColumnExpr("category_id, MAX(\"order\") AS order").
		Group("category_id").Scan(ctx)
	if err != nil {
		return nil, err
	}

	var pairs []string
	for _, result := range processRulesSubquery {
		pair := fmt.Sprintf("('%s', %d)", result.CategoryID, result.Order)
		pairs = append(pairs, pair)
	}

	// Consulta principal que compara category_id e order com os valores da subconsulta
	err = tx.NewSelect().
		Model(&processRules).
		Where(fmt.Sprintf("(category_id, \"order\") IN (%s)", strings.Join(pairs, ","))).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	for _, processRule := range processRules {
		if processRule.ID == uuid.Nil {
			continue
		}

		mapProcesses[processRule.ID] = processRule.CategoryID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return mapProcesses, nil
}

func (r *ProcessRuleRepositoryBun) IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error) {
	mapProcessRules, err := r.GetMapProcessRulesByLastOrder(ctx)
	if err != nil {
		return false, err
	}

	if _, exists := mapProcessRules[id]; exists {
		return true, nil
	}

	return false, nil
}

func (r *ProcessRuleRepositoryBun) GetProcessRulesByCategoryId(ctx context.Context, id string) ([]model.ProcessRule, error) {
	processRules := []model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processRules).Where("category_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processRules, nil
}

func (r *ProcessRuleRepositoryBun) GetAllProcessRules(ctx context.Context, page, perPage int, isActive bool) ([]model.ProcessRule, int, error) {
	processRules := []model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	// Calculate offset
	offset := page * perPage

	// Get paginated process rules with filter, only from active categories
	if err := tx.NewSelect().
		Model(&processRules).
		Join("INNER JOIN product_categories AS pc ON pc.id = pr.category_id").
		Where("pr.is_active = ?", isActive).
		Where("pc.is_active = true").
		Order("pr.name ASC").
		Limit(perPage).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	// Get total count with same filters
	total, err := tx.NewSelect().
		Model(&model.ProcessRule{}).
		Join("INNER JOIN product_categories AS pc ON pc.id = pr.category_id").
		Where("pr.is_active = ?", isActive).
		Where("pc.is_active = true").
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return processRules, total, nil
}

func (r *ProcessRuleRepositoryBun) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]model.ProcessRule, error) {
	processRules := []model.ProcessRule{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&processRules).
		Relation("", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("status = ?", "started")
		}).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return processRules, nil
}
