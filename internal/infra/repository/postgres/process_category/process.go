package processRulerepositorybun

import (
	"context"
	"sync"

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
