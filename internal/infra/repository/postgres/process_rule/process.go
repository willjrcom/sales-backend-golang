package processrulerepositorybun

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProcessRuleCategoryRepositoryBun(db *bun.DB) *ProcessRuleCategoryRepositoryBun {
	return &ProcessRuleCategoryRepositoryBun{db: db}
}

func (r *ProcessRuleCategoryRepositoryBun) RegisterProcessRule(ctx context.Context, s *productentity.ProcessRule) error {
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

func (r *ProcessRuleCategoryRepositoryBun) UpdateProcessRule(ctx context.Context, s *productentity.ProcessRule) error {
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

func (r *ProcessRuleCategoryRepositoryBun) DeleteProcessRule(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&productentity.ProcessRule{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcessRuleCategoryRepositoryBun) GetProcessRuleById(ctx context.Context, id string) (*productentity.ProcessRule, error) {
	processRule := &productentity.ProcessRule{}

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

func (r *ProcessRuleCategoryRepositoryBun) GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*productentity.ProcessRule, error) {
	processRule := &productentity.ProcessRule{}

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

func (r *ProcessRuleCategoryRepositoryBun) GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRules := []productentity.ProcessRule{}
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

func (r *ProcessRuleCategoryRepositoryBun) GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	processRulesSubquery := []productentity.ProcessRule{}
	processRules := []productentity.ProcessRule{}
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

	for _, processRule := range processRulesSubquery {
		mapProcesses[processRule.CategoryID] = processRule.ID
	}

	return mapProcesses, nil
}

func (r *ProcessRuleCategoryRepositoryBun) IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error) {
	mapProcessRules, err := r.GetMapProcessRulesByLastOrder(ctx)
	if err != nil {
		return false, err
	}

	return mapProcessRules[id] == uuid.Nil, nil
}
