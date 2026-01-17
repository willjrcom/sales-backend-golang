package model

import (
	"context"

	"github.com/google/uuid"
)

type ProcessRuleRepository interface {
	CreateProcessRule(ctx context.Context, ProcessRule *ProcessRule) error
	UpdateProcessRule(ctx context.Context, ProcessRule *ProcessRule) error
	DeleteProcessRule(ctx context.Context, id string) error
	GetProcessRuleById(ctx context.Context, id string) (*ProcessRule, error)
	GetProcessRulesByCategoryId(ctx context.Context, id string) ([]ProcessRule, error)
	GetAllProcessRules(ctx context.Context, page, perPage int, isActive bool) ([]ProcessRule, int, error)
	GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]ProcessRule, error)
	GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*ProcessRule, error)
	GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*ProcessRule, error)
	GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error)
}
