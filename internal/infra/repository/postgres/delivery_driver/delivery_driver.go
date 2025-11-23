package deliverydriverrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type DeliveryDriverRepositoryBun struct {
	db *bun.DB
}

func NewDeliveryDriverRepositoryBun(db *bun.DB) model.DeliveryDriverRepository {
	return &DeliveryDriverRepositoryBun{db: db}
}

func (r *DeliveryDriverRepositoryBun) CreateDeliveryDriver(ctx context.Context, s *model.DeliveryDriver) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *DeliveryDriverRepositoryBun) UpdateDeliveryDriver(ctx context.Context, s *model.DeliveryDriver) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *DeliveryDriverRepositoryBun) DeleteDeliveryDriver(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.DeliveryDriver{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *DeliveryDriverRepositoryBun) GetDeliveryDriverById(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	deliveryDriver := &model.DeliveryDriver{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(deliveryDriver).Where("driver.id = ?", id).Relation("Employee.User").Relation("OrderDeliveries").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deliveryDriver, nil
}

func (r *DeliveryDriverRepositoryBun) GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	deliveryDriver := &model.DeliveryDriver{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(deliveryDriver).Where("driver.employee_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deliveryDriver, nil
}

func (r *DeliveryDriverRepositoryBun) GetAllDeliveryDrivers(ctx context.Context) ([]model.DeliveryDriver, error) {
	deliveryDrivers := []model.DeliveryDriver{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&deliveryDrivers).Relation("Employee.User").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deliveryDrivers, nil
}
