package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type PickupOrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewPickupOrderRepositoryBun(db *bun.DB) *PickupOrderRepositoryBun {
	return &PickupOrderRepositoryBun{db: db}
}

func (r *PickupOrderRepositoryBun) CreatePickupOrder(ctx context.Context, pickupOrder *orderentity.PickupOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(pickupOrder).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PickupOrderRepositoryBun) UpdatePickupOrder(ctx context.Context, pickupOrder *orderentity.PickupOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(pickupOrder).WherePK().Where("id = ?", pickupOrder.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PickupOrderRepositoryBun) DeletePickupOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&orderentity.PickupOrder{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PickupOrderRepositoryBun) GetAllPickups(ctx context.Context) ([]orderentity.PickupOrder, error) {
	pickups := []orderentity.PickupOrder{}

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

func (r *PickupOrderRepositoryBun) GetPickupById(ctx context.Context, id string) (*orderentity.PickupOrder, error) {
	pickupOrder := &orderentity.PickupOrder{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(pickupOrder).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return pickupOrder, nil
}
