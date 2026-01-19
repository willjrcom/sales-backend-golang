package sizerepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SizeRepositoryBun struct {
	db *bun.DB
}

func NewSizeRepositoryBun(db *bun.DB) model.SizeRepository {
	return &SizeRepositoryBun{db: db}
}

func (r *SizeRepositoryBun) CreateSize(ctx context.Context, s *model.Size) error {

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

func (r *SizeRepositoryBun) UpdateSize(ctx context.Context, s *model.Size) error {

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

func (r *SizeRepositoryBun) DeleteSize(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false
	isActive := false
	if _, err := tx.NewUpdate().
		Model(&model.Size{}).
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

func (r *SizeRepositoryBun) GetSizeById(ctx context.Context, id string) (*model.Size, error) {
	size := &model.Size{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(size).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return size, nil
}

func (r *SizeRepositoryBun) GetSizeByIdWithProducts(ctx context.Context, id string) (*model.Size, error) {
	size := &model.Size{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(size).Where("id = ?", id).Relation("Products").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return size, nil
}

func (r *SizeRepositoryBun) GetSizesByCategoryId(ctx context.Context, categoryId string) ([]*model.Size, error) {
	sizes := []*model.Size{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&sizes).Where("category_id = ?", categoryId).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sizes, nil
}
