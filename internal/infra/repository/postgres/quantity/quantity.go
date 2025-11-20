package quantityrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QuantityRepositoryBun struct {
	db *bun.DB
}

func NewQuantityRepositoryBun(db *bun.DB) model.QuantityRepository {
	return &QuantityRepositoryBun{db: db}
}

func (r *QuantityRepositoryBun) CreateQuantity(ctx context.Context, s *model.Quantity) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QuantityRepositoryBun) UpdateQuantity(ctx context.Context, s *model.Quantity) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QuantityRepositoryBun) DeleteQuantity(ctx context.Context, id string) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Quantity{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QuantityRepositoryBun) GetQuantityById(ctx context.Context, id string) (*model.Quantity, error) {
	quantity := &model.Quantity{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(quantity).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return quantity, nil
}
