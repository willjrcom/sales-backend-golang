package sizerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewSizeCategoryRepositoryBun(ctx context.Context, db *bun.DB) *SizeCategoryRepositoryBun {
	return &SizeCategoryRepositoryBun{db: db}
}

func (r *SizeCategoryRepositoryBun) RegisterSize(ctx context.Context, s *productentity.Size) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *SizeCategoryRepositoryBun) UpdateSize(ctx context.Context, s *productentity.Size) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *SizeCategoryRepositoryBun) DeleteSize(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&productentity.Size{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *SizeCategoryRepositoryBun) GetSizeById(ctx context.Context, id string) (*productentity.Size, error) {
	size := &productentity.Size{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(size).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return size, nil
}
