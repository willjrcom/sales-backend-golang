package orderrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderPickupRepositoryBun struct {
	db *bun.DB
}

func NewOrderPickupRepositoryBun(db *bun.DB) model.OrderPickupRepository {
	return &OrderPickupRepositoryBun{db: db}
}

func (r *OrderPickupRepositoryBun) CreateOrderPickup(ctx context.Context, orderPickup *model.OrderPickup) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(orderPickup).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderPickupRepositoryBun) UpdateOrderPickup(ctx context.Context, orderPickup *model.OrderPickup) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(orderPickup).WherePK().Where("id = ?", orderPickup.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderPickupRepositoryBun) DeleteOrderPickup(ctx context.Context, id string) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderPickup{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderPickupRepositoryBun) GetAllPickups(ctx context.Context) ([]model.OrderPickup, error) {
	pickups := []model.OrderPickup{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&pickups).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pickups, nil
}

func (r *OrderPickupRepositoryBun) GetPickupById(ctx context.Context, id string) (*model.OrderPickup, error) {
	orderPickup := &model.OrderPickup{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(orderPickup).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return orderPickup, nil
}
