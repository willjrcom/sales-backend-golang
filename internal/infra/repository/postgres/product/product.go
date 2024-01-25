package productrepositorybun

import (
	"context"
	"sync"

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
	_, err := r.db.NewDelete().Model(&productentity.Product{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) GetProductById(ctx context.Context, id string) (*productentity.Product, error) {
	product := &productentity.Product{}

	ChangeSchema(r.db, "patrik_dog")

	r.mu.Lock()
	err := r.db.NewSelect().Model(product).Where("product.id = ?", id).Relation("Category").Relation("Size").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepositoryBun) GetAllProducts(ctx context.Context) ([]productentity.Product, error) {
	products := []productentity.Product{}

	ChangeSchema(r.db, "public")

	r.mu.Lock()
	err := r.db.NewSelect().Model(&products).Relation("Category").Relation("Size").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Exemplo de alteração dinâmica do schema no Bun ORM
func ChangeSchema(db *bun.DB, schemaName string) error {
	_, err := db.Exec("SET search_path=?", schemaName)
	return err
}
