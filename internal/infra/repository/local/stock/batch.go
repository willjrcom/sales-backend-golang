package stocklocal

import (
	"context"
	"errors"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// StockBatchRepositoryLocal is an in-memory implementation of model.StockBatchRepository.
type StockBatchRepositoryLocal struct {
	mu      sync.RWMutex
	batches map[string]*model.StockBatch
}

func NewStockBatchRepositoryLocal() *StockBatchRepositoryLocal {
	return &StockBatchRepositoryLocal{batches: make(map[string]*model.StockBatch)}
}

func (r *StockBatchRepositoryLocal) CreateBatch(ctx context.Context, db bun.IDB, b *model.StockBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.batches[b.ID.String()] = b
	return nil
}

func (r *StockBatchRepositoryLocal) UpdateBatch(ctx context.Context, db bun.IDB, b *model.StockBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.batches[b.ID.String()] = b
	return nil
}

func (r *StockBatchRepositoryLocal) GetBatchByID(ctx context.Context, db bun.IDB, id string) (*model.StockBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.batches[id]
	if !ok {
		return nil, errors.New("batch not found: " + id)
	}
	return b, nil
}

func (r *StockBatchRepositoryLocal) GetActiveBatchesByStockID(ctx context.Context, stockID string) ([]model.StockBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockBatch
	for _, b := range r.batches {
		if b.StockID.String() != stockID {
			continue
		}
		batch := b.ToDomain()
		if batch.CurrentQuantity.IsPositive() && !batch.IsExpired() {
			result = append(result, *b)
		}
	}
	// Sort by created_at ascending (FIFO) — simplification: order by ID string
	return result, nil
}

func (r *StockBatchRepositoryLocal) GetBatchesByStockID(ctx context.Context, stockID string) ([]model.StockBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.StockBatch
	for _, b := range r.batches {
		if b.StockID.String() != stockID {
			continue
		}
		batch := b.ToDomain()
		if batch.CurrentQuantity.IsPositive() {
			result = append(result, *b)
		}
	}
	return result, nil
}

func (r *StockBatchRepositoryLocal) GetActiveBatchesByStockIDForUpdate(ctx context.Context, db bun.IDB, stockID string) ([]model.StockBatch, error) {
	return r.GetActiveBatchesByStockID(ctx, stockID)
}
