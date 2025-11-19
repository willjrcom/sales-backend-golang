package orderrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderDeliveryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderDeliveryRepositoryBun(db *bun.DB) model.OrderDeliveryRepository {
	return &OrderDeliveryRepositoryBun{db: db}
}

func (r *OrderDeliveryRepositoryBun) CreateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(delivery).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderDeliveryRepositoryBun) UpdateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(delivery).Where("id = ?", delivery.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderDeliveryRepositoryBun) DeleteOrderDelivery(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderDelivery{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderDeliveryRepositoryBun) GetAllDeliveries(ctx context.Context) ([]model.OrderDelivery, error) {
	deliveries := []model.OrderDelivery{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&deliveries).Where("delivery.status != ?", orderentity.OrderDeliveryStatusStaging).Relation("Client").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *OrderDeliveryRepositoryBun) GetDeliveryById(ctx context.Context, id string) (*model.OrderDelivery, error) {
	delivery := &model.OrderDelivery{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(delivery).Where("delivery.id = ?", id).Relation("Client.Address").Relation("Address").Relation("Driver").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return delivery, nil
}

func (r *OrderDeliveryRepositoryBun) GetDeliveriesByIds(ctx context.Context, ids []string) ([]model.OrderDelivery, error) {
	deliveries := []model.OrderDelivery{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&deliveries).Where("delivery.id IN (?)", bun.In(ids)).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deliveries, nil
}
