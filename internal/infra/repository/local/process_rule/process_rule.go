package processrulerepositorylocal

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProcessRuleRepositoryLocal struct {
	mu    sync.RWMutex
	rules map[string]*model.ProcessRule
}

func NewProcessRuleRepositoryLocal() model.ProcessRuleRepository {
	return &ProcessRuleRepositoryLocal{rules: make(map[string]*model.ProcessRule)}
}

func (r *ProcessRuleRepositoryLocal) CreateProcessRule(ctx context.Context, pr *model.ProcessRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[pr.ID.String()] = pr
	return nil
}

func (r *ProcessRuleRepositoryLocal) UpdateProcessRule(ctx context.Context, pr *model.ProcessRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[pr.ID.String()] = pr
	return nil
}

func (r *ProcessRuleRepositoryLocal) DeleteProcessRule(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}

func (r *ProcessRuleRepositoryLocal) GetProcessRuleById(ctx context.Context, id string) (*model.ProcessRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if pr, ok := r.rules[id]; ok {
		return pr, nil
	}
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetProcessRulesByCategoryId(ctx context.Context, id string) ([]model.ProcessRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []model.ProcessRule{}
	for _, pr := range r.rules {
		if pr.CategoryID.String() == id {
			out = append(out, *pr)
		}
	}
	return out, nil
}

func (r *ProcessRuleRepositoryLocal) GetAllProcessRules(ctx context.Context, page, perPage int, isActive bool) ([]model.ProcessRule, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.ProcessRule, 0, len(r.rules))
	for _, pr := range r.rules {
		if pr.IsActive == isActive {
			out = append(out, *pr)
		}
	}
	total := len(out)

	// Apply pagination
	start := page * perPage
	end := start + perPage
	if start > total {
		return []model.ProcessRule{}, total, nil
	}
	if end > total {
		end = total
	}

	return out[start:end], total, nil
}

func (r *ProcessRuleRepositoryLocal) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]model.ProcessRule, error) {
	rules, _, err := r.GetAllProcessRules(ctx, 0, 1000, true)
	return rules, err
}

func (r *ProcessRuleRepositoryLocal) GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*model.ProcessRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, pr := range r.rules {
		if pr.CategoryID.String() == id && pr.Order == order {
			return pr, nil
		}
	}
	return nil, nil
}

func (r *ProcessRuleRepositoryLocal) GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*model.ProcessRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var first *model.ProcessRule
	for _, pr := range r.rules {
		if pr.CategoryID.String() == id {
			if first == nil || pr.Order < first.Order {
				first = pr
			}
		}
	}
	return first, nil
}

func (r *ProcessRuleRepositoryLocal) GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := make(map[uuid.UUID]uuid.UUID)
	for _, pr := range r.rules {
		if first, _ := r.GetFirstProcessRuleByCategoryId(ctx, pr.CategoryID.String()); first != nil {
			m[pr.CategoryID] = first.ID
		}
	}
	return m, nil
}

func (r *ProcessRuleRepositoryLocal) GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	temp := map[uuid.UUID]*model.ProcessRule{}
	for _, pr := range r.rules {
		cid := pr.CategoryID
		if cur, ok := temp[cid]; !ok || pr.Order > cur.Order {
			temp[cid] = pr
		}
	}
	out := make(map[uuid.UUID]uuid.UUID, len(temp))
	for cid, pr := range temp {
		out[cid] = pr.ID
	}
	return out, nil
}

func (r *ProcessRuleRepositoryLocal) IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pr, ok := r.rules[id.String()]
	if !ok {
		return false, nil
	}
	lastMap, _ := r.GetMapProcessRulesByLastOrder(ctx)
	if lastID, ok2 := lastMap[pr.CategoryID]; ok2 {
		return lastID == id, nil
	}
	return false, nil
}
