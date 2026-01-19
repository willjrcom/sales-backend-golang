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

func (r *QuantityRepositoryBun) UpdateQuantity(ctx context.Context, s *model.Quantity) error {

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

func (r *QuantityRepositoryBun) DeleteQuantity(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false
	isActive := false
	if _, err := tx.NewUpdate().
		Model(&model.Quantity{}).
		Set("is_active = ?", isActive).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QuantityRepositoryBun) GetQuantityById(ctx context.Context, id string) (*model.Quantity, error) {
	quantity := &model.Quantity{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(quantity).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return quantity, nil
}

func (r *QuantityRepositoryBun) GetQuantitiesByCategoryId(ctx context.Context, categoryId string) ([]*model.Quantity, error) {
	quantities := []*model.Quantity{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&quantities).Where("category_id = ?", categoryId).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return quantities, nil
}
