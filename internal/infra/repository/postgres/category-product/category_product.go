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

func NewCategoryProductRepositoryBun(ctx context.Context, db *bun.DB) *CategoryProductRepositoryBun {
	return &CategoryProductRepositoryBun{db: db}
}

func (r *CategoryProductRepositoryBun) RegisterCategoryProduct(ctx context.Context, cp *productentity.CategoryProduct) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(cp).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryProductRepositoryBun) UpdateCategoryProduct(ctx context.Context, cp *productentity.CategoryProduct) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(cp).Where("id = ?", cp.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryProductRepositoryBun) DeleteCategoryProduct(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&productentity.CategoryProduct{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *CategoryProductRepositoryBun) GetCategoryProductById(ctx context.Context, id string) (*productentity.CategoryProduct, error) {
	category := &productentity.CategoryProduct{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryProductRepositoryBun) GetAllCategoryProduct(ctx context.Context) ([]productentity.CategoryProduct, error) {
	categories := []productentity.CategoryProduct{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(&categories).Relation("Products").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return categories, nil
}
