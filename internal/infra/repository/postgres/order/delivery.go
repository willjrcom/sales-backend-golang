package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type DeliveryOrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewDeliveryOrderRepositoryBun(db *bun.DB) *DeliveryOrderRepositoryBun {
	return &DeliveryOrderRepositoryBun{db: db}
}

func (r *DeliveryOrderRepositoryBun) CreateDeliveryOrder(ctx context.Context, delivery *orderentity.DeliveryOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(delivery).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryOrderRepositoryBun) UpdateDeliveryOrder(ctx context.Context, delivery *orderentity.DeliveryOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(delivery).WherePK().Where("id = ?", delivery.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryOrderRepositoryBun) DeleteDeliveryOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&orderentity.DeliveryOrder{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryOrderRepositoryBun) GetAllDeliveries(ctx context.Context) ([]orderentity.DeliveryOrder, error) {
	deliveries := []orderentity.DeliveryOrder{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&deliveries).Relation("Client").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryOrderRepositoryBun) GetDeliveryById(ctx context.Context, id string) (*orderentity.DeliveryOrder, error) {
	delivery := &orderentity.DeliveryOrder{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(delivery).Where("id = ?", id).Relation("Client").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	return delivery, nil
}
