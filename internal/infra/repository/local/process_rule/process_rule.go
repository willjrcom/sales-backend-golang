package processrulerepositorylocal

import (
	"context"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProcessRuleRepositoryLocal struct {}

func NewProcessRuleRepositoryLocal() model.ProcessRuleRepository {
	return &ProcessRuleRepositoryLocal{}
}

func (r *ProcessRuleRepositoryLocal) CreateProcessRule(ctx context.Context, ProcessRule *model.ProcessRule) error {
	return nil
}

func (r *ProcessRuleRepositoryLocal) UpdateProcessRule(ctx context.Context, ProcessRule *model.ProcessRule) error {
	return nil
}

func (r *ProcessRuleRepositoryLocal) DeleteProcessRule(ctx context.Context, id string) error {
	return nil
}

func (r *ProcessRuleRepositoryLocal) GetProcessRuleById(ctx context.Context, id string) (*model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetProcessRulesByCategoryId(ctx context.Context, id string) ([]model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetAllProcessRules(ctx context.Context) ([]model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*model.ProcessRule, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}
