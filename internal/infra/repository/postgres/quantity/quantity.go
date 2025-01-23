package quantityrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QuantityRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewQuantityRepositoryBun(db *bun.DB) model.QuantityRepository {
	return &QuantityRepositoryBun{db: db}
}

func (r *QuantityRepositoryBun) CreateQuantity(ctx context.Context, s *model.Quantity) error {
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

func (r *QuantityRepositoryBun) UpdateQuantity(ctx context.Context, s *model.Quantity) error {
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

func (r *QuantityRepositoryBun) DeleteQuantity(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.Quantity{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *QuantityRepositoryBun) GetQuantityById(ctx context.Context, id string) (*model.Quantity, error) {
	quantity := &model.Quantity{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(quantity).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return quantity, nil
}
