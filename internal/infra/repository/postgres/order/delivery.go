package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
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
	_, err := r.db.NewInsert().Model(delivery).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *DeliveryOrderRepositoryBun) UpdateDeliveryOrder(ctx context.Context, delivery *orderentity.DeliveryOrder) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(delivery).WherePK().Where("id = ?", delivery.ID).Exec(ctx)
	r.mu.Unlock()

	return err
}

func (r *DeliveryOrderRepositoryBun) DeleteDeliveryOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&orderentity.DeliveryOrder{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *DeliveryOrderRepositoryBun) GetAllDeliveries(ctx context.Context) ([]orderentity.DeliveryOrder, error) {
	deliveries := []orderentity.DeliveryOrder{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&deliveries).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryOrderRepositoryBun) GetDeliveryById(ctx context.Context, id string) (*orderentity.DeliveryOrder, error) {
	delivery := &orderentity.DeliveryOrder{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(delivery).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return delivery, nil
}
