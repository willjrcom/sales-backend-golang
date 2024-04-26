package categoryrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryProductRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewCategoryProductRepositoryBun(db *bun.DB) *CategoryProductRepositoryBun {
	return &CategoryProductRepositoryBun{db: db}
}

func (r *CategoryProductRepositoryBun) RegisterCategory(ctx context.Context, cp *productentity.Category) error {
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
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	return r.updateAdditionalCategories(ctx, tx, cp.ID, cp.AdditionalCategories)
}

func (r *CategoryProductRepositoryBun) UpdateCategory(ctx context.Context, c *productentity.Category) error {
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
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	return r.updateAdditionalCategories(ctx, tx, c.ID, c.AdditionalCategories)
}

func (r *CategoryProductRepositoryBun) updateAdditionalCategories(ctx context.Context, tx bun.Tx, categoryID uuid.UUID, additionalCategories []productentity.Category) error {
	if err := database.ChangeSchema(ctx, r.db); err != nil {

		return err
	}

	if _, err := tx.NewDelete().Model(&productentity.CategoryToAdditional{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	for _, ac := range additionalCategories {
		categoryToAdditional := &productentity.CategoryToAdditional{
			CategoryID:           categoryID,
			AdditionalCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToAdditional).Exec(ctx); err != nil {
			if errRollBack := tx.Rollback(); errRollBack != nil {
				return errRollBack
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CategoryProductRepositoryBun) DeleteCategory(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&productentity.Category{}).Where("id = ?", id).Exec(ctx); err != nil {
		if errRoolback := tx.Rollback(); errRoolback != nil {
			return errRoolback
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&productentity.CategoryToAdditional{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&productentity.CategoryToAdditional{}).Where("additional_category_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	return nil
}

func (r *CategoryProductRepositoryBun) GetCategoryById(ctx context.Context, id string) (*productentity.Category, error) {
	category := &productentity.Category{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryProductRepositoryBun) GetCategoryByName(ctx context.Context, name string, withRelation bool) (*productentity.Category, error) {
	category := &productentity.Category{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(category).Where("name = ?", name)

	if withRelation {
		query.Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryProductRepositoryBun) GetAllCategories(ctx context.Context) ([]productentity.Category, error) {
	categories := []productentity.Category{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&categories).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Scan(ctx); err != nil {
		return nil, err
	}

	return categories, nil
}
