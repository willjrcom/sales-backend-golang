package companycategoryrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyCategoryRepositoryBun struct {
	db *bun.DB
}

func NewCompanyCategoryRepositoryBun(db *bun.DB) model.CompanyCategoryRepository {
	return &CompanyCategoryRepositoryBun{db: db}
}

func (r *CompanyCategoryRepositoryBun) Create(ctx context.Context, category *model.CompanyCategory) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err = tx.NewInsert().Model(category).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyCategoryRepositoryBun) Update(ctx context.Context, category *model.CompanyCategory) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err = tx.NewUpdate().Model(category).Where("id = ?", category.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyCategoryRepositoryBun) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.CompanyCategory{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyCategoryRepositoryBun) GetByID(ctx context.Context, id uuid.UUID) (*model.CompanyCategory, error) {
	category := &model.CompanyCategory{}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(category).Relation("Sponsors").Relation("Advertisements").Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CompanyCategoryRepositoryBun) GetAllCompanyCategories(ctx context.Context) ([]model.CompanyCategory, error) {
	categories := []model.CompanyCategory{}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&categories).Relation("Sponsors").Relation("Advertisements").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return categories, nil
}
