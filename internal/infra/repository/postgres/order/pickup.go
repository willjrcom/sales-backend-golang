package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderPickupRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderPickupRepositoryBun(db *bun.DB) model.OrderPickupRepository {
	return &OrderPickupRepositoryBun{db: db}
}

func (r *OrderPickupRepositoryBun) CreateOrderPickup(ctx context.Context, orderPickup *model.OrderPickup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(orderPickup).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderPickupRepositoryBun) UpdateOrderPickup(ctx context.Context, orderPickup *model.OrderPickup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(orderPickup).WherePK().Where("id = ?", orderPickup.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderPickupRepositoryBun) DeleteOrderPickup(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.OrderPickup{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderPickupRepositoryBun) GetAllPickups(ctx context.Context) ([]model.OrderPickup, error) {
	pickups := []model.OrderPickup{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&pickups).Scan(ctx); err != nil {
		return nil, err
	}

	return pickups, nil
}

func (r *OrderPickupRepositoryBun) GetPickupById(ctx context.Context, id string) (*model.OrderPickup, error) {
	orderPickup := &model.OrderPickup{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(orderPickup).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return orderPickup, nil
}
