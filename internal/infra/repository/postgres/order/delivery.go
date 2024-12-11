package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderDeliveryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderDeliveryRepositoryBun(db *bun.DB) *OrderDeliveryRepositoryBun {
	return &OrderDeliveryRepositoryBun{db: db}
}

func (r *OrderDeliveryRepositoryBun) CreateOrderDelivery(ctx context.Context, delivery *orderentity.OrderDelivery) error {
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

func (r *OrderDeliveryRepositoryBun) UpdateOrderDelivery(ctx context.Context, delivery *orderentity.OrderDelivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(delivery).Where("id = ?", delivery.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderDeliveryRepositoryBun) DeleteOrderDelivery(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&orderentity.OrderDelivery{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderDeliveryRepositoryBun) GetAllDeliveries(ctx context.Context) ([]orderentity.OrderDelivery, error) {
	deliveries := []orderentity.OrderDelivery{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&deliveries).Where("delivery.status != ?", orderentity.OrderDeliveryStatusStaging).Relation("Client").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *OrderDeliveryRepositoryBun) GetDeliveryById(ctx context.Context, id string) (*orderentity.OrderDelivery, error) {
	delivery := &orderentity.OrderDelivery{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(delivery).Where("delivery.id = ?", id).Relation("Client.Address").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	return delivery, nil
}
