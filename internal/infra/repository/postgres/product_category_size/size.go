package productcategorysizerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type SizeRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewSizeRepositoryBun(db *bun.DB) *SizeRepositoryBun {
	return &SizeRepositoryBun{db: db}
}

func (r *SizeRepositoryBun) RegisterSize(ctx context.Context, s *productentity.Size) error {
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

func (r *SizeRepositoryBun) UpdateSize(ctx context.Context, s *productentity.Size) error {
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

func (r *SizeRepositoryBun) DeleteSize(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&productentity.Size{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *SizeRepositoryBun) GetSizeById(ctx context.Context, id string) (*productentity.Size, error) {
	size := &productentity.Size{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(size).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return size, nil
}
