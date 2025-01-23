package processrulerepositorybun

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProcessRuleRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProcessRuleRepositoryBun(db *bun.DB) model.ProcessRuleRepository {
	return &ProcessRuleRepositoryBun{db: db}
}

func (r *ProcessRuleRepositoryBun) CreateProcessRule(ctx context.Context, s *model.ProcessRule) error {
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

func (r *ProcessRuleRepositoryBun) UpdateProcessRule(ctx context.Context, s *model.ProcessRule) error {
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

func (r *ProcessRuleRepositoryBun) DeleteProcessRule(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.ProcessRule{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcessRuleRepositoryBun) GetProcessRuleById(ctx context.Context, id string) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(processRule).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(processRule).Where("category_id = ? and order = 1", id).Scan(ctx); err != nil {
		return nil, err
	}

	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRules := []model.ProcessRule{}
	mapProcesses := map[uuid.UUID]uuid.UUID{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processRules).Where("\"order\" = 1").Scan(ctx); err != nil {
		return nil, err
	}

	for _, processRule := range processRules {
		mapProcesses[processRule.CategoryID] = processRule.ID
	}

	return mapProcesses, nil
}

func (r *ProcessRuleRepositoryBun) GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*model.ProcessRule, error) {
	processRule := &model.ProcessRule{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(processRule).Where("category_id = ? and \"order\" = ?", id, order).Scan(ctx); err != nil {
		return nil, err
	}

	return processRule, nil
}

func (r *ProcessRuleRepositoryBun) GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRulesSubquery := []model.ProcessRule{}
	processRules := []model.ProcessRule{}
	mapProcesses := map[uuid.UUID]uuid.UUID{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	// Subconsulta para obter o máximo order para cada category_id
	err := r.db.NewSelect().Model(&processRulesSubquery).ColumnExpr("category_id, MAX(\"order\") AS order").
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
	err = r.db.NewSelect().
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

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processRules).Where("category_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return processRules, nil
}

func (r *ProcessRuleRepositoryBun) GetAllProcessRules(ctx context.Context) ([]model.ProcessRule, error) {
	processRules := []model.ProcessRule{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processRules).Scan(ctx); err != nil {
		return nil, err
	}

	return processRules, nil
}

func (r *ProcessRuleRepositoryBun) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]model.ProcessRule, error) {
	processRules := []model.ProcessRule{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&processRules).
		Relation("", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("status = ?", "started")
		}).Scan(ctx); err != nil {
		return nil, err
	}

	return processRules, nil
}
