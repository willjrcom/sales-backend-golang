package stocklocal

import (
	"context"
	"errors"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// StockAlertRepositoryLocal is an in-memory implementation of model.StockAlertRepository.
type StockAlertRepositoryLocal struct {
	mu     sync.RWMutex
	alerts map[string]*model.StockAlert
}

func NewStockAlertRepositoryLocal() *StockAlertRepositoryLocal {
	return &StockAlertRepositoryLocal{alerts: make(map[string]*model.StockAlert)}
}

func (r *StockAlertRepositoryLocal) CreateAlert(ctx context.Context, a *model.StockAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts[a.ID.String()] = a
	return nil
}

func (r *StockAlertRepositoryLocal) UpdateAlert(ctx context.Context, a *model.StockAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts[a.ID.String()] = a
	return nil
}

func (r *StockAlertRepositoryLocal) DeleteAlert(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.alerts, id)
	return nil
}

func (r *StockAlertRepositoryLocal) GetAlertByID(ctx context.Context, id string) (*model.StockAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.alerts[id]
	if !ok {
		return nil, errors.New("alert not found: " + id)
	}
	return a, nil
}

func (r *StockAlertRepositoryLocal) GetActiveAlerts(ctx context.Context) ([]model.StockAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockAlert
	for _, a := range r.alerts {
		if !a.IsResolved {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (r *StockAlertRepositoryLocal) GetAllAlerts(ctx context.Context) ([]model.StockAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockAlert
	for _, a := range r.alerts {
		result = append(result, *a)
	}
	return result, nil
}

func (r *StockAlertRepositoryLocal) GetAlertsByStockID(ctx context.Context, stockID string) ([]model.StockAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockAlert
	for _, a := range r.alerts {
		if a.StockID.String() == stockID {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (r *StockAlertRepositoryLocal) GetResolvedAlerts(ctx context.Context) ([]model.StockAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockAlert
	for _, a := range r.alerts {
		if a.IsResolved {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (r *StockAlertRepositoryLocal) ResolveAlert(ctx context.Context, alertID string, resolvedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.alerts[alertID]
	if !ok {
		return errors.New("alert not found: " + alertID)
	}
	a.IsResolved = true
	return nil
}
