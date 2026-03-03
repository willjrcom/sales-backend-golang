package stocklocal

import (
	"context"
	"errors"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockMovementRepositoryLocal struct {
	mu        sync.RWMutex
	movements map[string]*model.StockMovement
}

func NewStockMovementRepositoryLocal() *StockMovementRepositoryLocal {
	return &StockMovementRepositoryLocal{movements: make(map[string]*model.StockMovement)}
}

func (r *StockMovementRepositoryLocal) CreateMovement(ctx context.Context, db bun.IDB, m *model.StockMovement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.movements[m.ID.String()] = m
	return nil
}

func (r *StockMovementRepositoryLocal) GetMovementsByStockID(ctx context.Context, stockID string, date *string) ([]model.StockMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockMovement
	for _, m := range r.movements {
		if m.StockID.String() == stockID {
			result = append(result, *m)
		}
	}
	return result, nil
}

func (r *StockMovementRepositoryLocal) GetMovementsByProductID(ctx context.Context, productID string) ([]model.StockMovement, error) {
	return nil, nil
}

func (r *StockMovementRepositoryLocal) GetMovementsByOrderID(ctx context.Context, orderID string) ([]model.StockMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockMovement
	for _, m := range r.movements {
		if m.OrderID != nil && m.OrderID.String() == orderID {
			result = append(result, *m)
		}
	}
	return result, nil
}

func (r *StockMovementRepositoryLocal) GetAllMovements(ctx context.Context) ([]model.StockMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockMovement
	for _, m := range r.movements {
		result = append(result, *m)
	}
	return result, nil
}

func (r *StockMovementRepositoryLocal) GetMovementsByDateRange(ctx context.Context, start, end string) ([]model.StockMovement, error) {
	return nil, errors.New("not implemented in local repo")
}
