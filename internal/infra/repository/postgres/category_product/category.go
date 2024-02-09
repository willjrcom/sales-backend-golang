package categoryrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
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
	_, err := r.db.NewInsert().Model(cp).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryProductRepositoryBun) UpdateCategory(ctx context.Context, c *productentity.Category) error {
	r.mu.Lock()
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		r.mu.Unlock()
		return err
	}

	if _, err = tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		if err = tx.Rollback(); err != nil {
			r.mu.Unlock()
			return err
		}

		r.mu.Unlock()
		return err
	}

	if _, err = tx.NewDelete().Model(productentity.CategoryCategory{}).Where("category_id = ?", c.ID).Exec(ctx); err != nil {
		if err = tx.Rollback(); err != nil {
			r.mu.Unlock()
			return err
		}

		r.mu.Unlock()
		return err
	}

	for _, ac := range c.AdditionalCategories {
		if _, err = tx.NewInsert().Model(ac).Exec(ctx); err != nil {
			if err = tx.Rollback(); err != nil {
				r.mu.Unlock()
				return err
			}

			r.mu.Unlock()
			return err
		}
	}

	r.mu.Unlock()

	return nil
}

func (r *CategoryProductRepositoryBun) DeleteCategory(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&productentity.Category{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *CategoryProductRepositoryBun) GetCategoryById(ctx context.Context, id string) (*productentity.Category, error) {
	category := &productentity.Category{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("Processes").Relation("AditionalCategories").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryProductRepositoryBun) GetAllCategories(ctx context.Context) ([]productentity.Category, error) {
	categories := []productentity.Category{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(&categories).Relation("Sizes").Relation("Quantities").Relation("Processes").Relation("AditionalCategories").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return categories, nil
}
