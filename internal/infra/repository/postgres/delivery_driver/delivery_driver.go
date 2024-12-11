package deliverydriverrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type DeliveryDriverRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewDeliveryDriverRepositoryBun(db *bun.DB) *DeliveryDriverRepositoryBun {
	return &DeliveryDriverRepositoryBun{db: db}
}

func (r *DeliveryDriverRepositoryBun) CreateDeliveryDriver(ctx context.Context, s *orderentity.DeliveryDriver) error {
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

func (r *DeliveryDriverRepositoryBun) UpdateDeliveryDriver(ctx context.Context, s *orderentity.DeliveryDriver) error {
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

func (r *DeliveryDriverRepositoryBun) DeleteDeliveryDriver(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&orderentity.DeliveryDriver{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryDriverRepositoryBun) GetDeliveryDriverById(ctx context.Context, id string) (*orderentity.DeliveryDriver, error) {
	deliveryDriver := &orderentity.DeliveryDriver{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(deliveryDriver).Where("driver.id = ?", id).Relation("Employee").Relation("OrderDeliveries").Scan(ctx); err != nil {
		return nil, err
	}

	return deliveryDriver, nil
}

func (r *DeliveryDriverRepositoryBun) GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*orderentity.DeliveryDriver, error) {
	deliveryDriver := &orderentity.DeliveryDriver{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(deliveryDriver).Where("driver.employee_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return deliveryDriver, nil
}

func (r *DeliveryDriverRepositoryBun) GetAllDeliveryDrivers(ctx context.Context) ([]orderentity.DeliveryDriver, error) {
	deliveryDrivers := []orderentity.DeliveryDriver{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&deliveryDrivers).Relation("Employee").Scan(ctx); err != nil {
		return nil, err
	}

	return deliveryDrivers, nil
}
