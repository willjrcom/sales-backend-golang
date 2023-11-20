package categoryrepositorybun

import (
	"context"
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

func (r *CategoryProductRepositoryBun) UpdateCategory(ctx context.Context, cp *productentity.Category) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(cp).Where("id = ?", cp.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

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
	err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryProductRepositoryBun) GetAllCategoryProducts(ctx context.Context) ([]productentity.Category, error) {
	categories := []productentity.Category{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(&categories).Relation("Products").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryProductRepositoryBun) GetAllCategorySizes(ctx context.Context) ([]productentity.Category, error) {
	categories := []productentity.Category{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(&categories).Relation("Sizes").Relation("Quantities").Relation("AditionalCategories").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return categories, nil
}
