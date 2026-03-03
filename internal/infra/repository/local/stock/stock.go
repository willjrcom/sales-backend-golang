package stocklocal

import (
	"context"
	"errors"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// StockRepositoryLocal is an in-memory implementation of model.StockRepository.
type StockRepositoryLocal struct {
	mu     sync.RWMutex
	stocks map[string]*model.Stock
}

func NewStockRepositoryLocal() *StockRepositoryLocal {
	return &StockRepositoryLocal{stocks: make(map[string]*model.Stock)}
}

func (r *StockRepositoryLocal) CreateStock(ctx context.Context, s *model.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stocks[s.ID.String()] = s
	return nil
}

func (r *StockRepositoryLocal) UpdateStock(ctx context.Context, db bun.IDB, s *model.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stocks[s.ID.String()] = s
	return nil
}

func (r *StockRepositoryLocal) GetStockByID(ctx context.Context, id string) (*model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.stocks[id]
	if !ok {
		return nil, errors.New("stock not found: " + id)
	}
	return s, nil
}

func (r *StockRepositoryLocal) GetStockByProductID(ctx context.Context, productID string) ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Stock
	for _, s := range r.stocks {
		if s.ProductID.String() == productID {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (r *StockRepositoryLocal) GetStockByVariationID(ctx context.Context, variationID string) (*model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.stocks {
		if s.ProductVariationID != nil && s.ProductVariationID.String() == variationID {
			return s, nil
		}
	}
	return nil, errors.New("stock not found for variation: " + variationID)
}

func (r *StockRepositoryLocal) GetAllStocks(ctx context.Context, page, perPage int) ([]model.Stock, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Stock
	for _, s := range r.stocks {
		result = append(result, *s)
	}
	return result, len(result), nil
}

func (r *StockRepositoryLocal) GetActiveStocks(ctx context.Context) ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Stock
	for _, s := range r.stocks {
		if s.IsActive {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (r *StockRepositoryLocal) GetLowStockProducts(ctx context.Context) ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Stock
	for _, s := range r.stocks {
		stock := s.ToDomain()
		if stock.IsLowStock() {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (r *StockRepositoryLocal) GetOutOfStockProducts(ctx context.Context) ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Stock
	for _, s := range r.stocks {
		stock := s.ToDomain()
		if stock.IsOutOfStock() {
			result = append(result, *s)
		}
	}
	return result, nil
}
