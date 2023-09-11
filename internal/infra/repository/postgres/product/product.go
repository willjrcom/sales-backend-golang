package productrepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProductRepositoryBun(db *bun.DB) *ProductRepositoryBun {
	return &ProductRepositoryBun{db: db}
}

func (r *ProductRepositoryBun) RegisterProduct(ctx context.Context, p *productentity.Product) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) UpdateProduct(ctx context.Context, p *productentity.Product) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) DeleteProduct(ctx context.Context, id string) error {
	r.mu.Lock()
	r.db.NewDelete().Model(&productentity.Product{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryBun) GetProductById(ctx context.Context, id string) (*productentity.Product, error) {
	product := &productentity.Product{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(product).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepositoryBun) GetProductsBy(ctx context.Context, p *productentity.Product) ([]productentity.Product, error) {
	products := []productentity.Product{}

	r.mu.Lock()
	query := r.db.NewSelect().Model(&productentity.Product{})

	if p.Code != "" {
		query.Where("code = ?", p.Code)
	}
	if p.Name != "" {
		query.Where("name = ?", p.Name)
	}
	if p.CategoryID != uuid.Nil {
		query.Where("categoryID = ?", p.Category.ID)
	}
	if p.Size != "" {
		query.Where("size = ?", p.Size)
	}

	err := query.Scan(ctx, &products)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepositoryBun) GetAllProductsByCategory(ctx context.Context, category string) ([]productentity.Product, error) {
	products := []productentity.Product{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&products).Where("category.name = ?", category).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return products, nil
}
