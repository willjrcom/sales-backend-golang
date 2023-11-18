package quantityrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type QuantityCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewQuantityCategoryRepositoryBun(db *bun.DB) *QuantityCategoryRepositoryBun {
	return &QuantityCategoryRepositoryBun{db: db}
}

func (r *QuantityCategoryRepositoryBun) RegisterQuantity(ctx context.Context, s *productentity.Quantity) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *QuantityCategoryRepositoryBun) UpdateQuantity(ctx context.Context, s *productentity.Quantity) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *QuantityCategoryRepositoryBun) DeleteQuantity(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&productentity.Quantity{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *QuantityCategoryRepositoryBun) GetQuantityById(ctx context.Context, id string) (*productentity.Quantity, error) {
	quantity := &productentity.Quantity{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(quantity).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return quantity, nil
}
