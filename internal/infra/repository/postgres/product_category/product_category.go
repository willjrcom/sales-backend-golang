package productcategoryrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProductCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProductCategoryRepositoryBun(db *bun.DB) *ProductCategoryRepositoryBun {
	return &ProductCategoryRepositoryBun{db: db}
}

func (r *ProductCategoryRepositoryBun) CreateCategory(ctx context.Context, cp *model.ProductCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(cp).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, &tx, cp.ID, cp.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, &tx, cp.ID, cp.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) UpdateCategory(ctx context.Context, c *model.ProductCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, &tx, c.ID, c.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, &tx, c.ID, c.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) updateAdditionalCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, additionalCategories []model.ProductCategory) error {
	if err := database.ChangeSchema(ctx, r.db); err != nil {

		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, ac := range additionalCategories {
		categoryToAdditional := &model.ProductCategoryToAdditional{
			CategoryID:           categoryID,
			AdditionalCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToAdditional).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) updateComplementCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, complementCategories []model.ProductCategory) error {
	if err := database.ChangeSchema(ctx, r.db); err != nil {

		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToComplement{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, ac := range complementCategories {
		categoryToComplement := &model.ProductCategoryToComplement{
			CategoryID:           categoryID,
			ComplementCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToComplement).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) DeleteCategory(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategory{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("additional_category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToComplement{}).Where("complement_category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Size{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Quantity{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProcessRule{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryById(ctx context.Context, id string) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories").Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryByName(ctx context.Context, name string, withRelation bool) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(category).Where("name = ?", name)

	if withRelation {
		query.Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&categories).Relation("Products.Size").Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories").Scan(ctx); err != nil {
		return nil, err
	}

	return categories, nil
}
