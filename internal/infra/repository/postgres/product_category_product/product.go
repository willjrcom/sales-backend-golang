package productcategoryproductrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProductRepositoryBun(db *bun.DB) *ProductRepositoryBun {
	return &ProductRepositoryBun{db: db}
}

func (r *ProductRepositoryBun) CreateProduct(ctx context.Context, p *productentity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) UpdateProduct(ctx context.Context, p *productentity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) DeleteProduct(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&productentity.Product{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) GetProductById(ctx context.Context, id string) (*productentity.Product, error) {
	product := &productentity.Product{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(product).Where("product.id = ?", id).Relation("Category").Relation("Size").Scan(ctx); err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepositoryBun) GetProductByCode(ctx context.Context, code string) (*productentity.Product, error) {
	product := &productentity.Product{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(product).Where("product.code = ?", code).Relation("Category").Relation("Size").Scan(ctx); err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepositoryBun) GetAllProducts(ctx context.Context) ([]productentity.Product, error) {
	products := []productentity.Product{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&products).Relation("Category").Relation("Size").Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}
