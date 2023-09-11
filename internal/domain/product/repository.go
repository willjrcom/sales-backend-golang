package productentity

import (
	"context"
)

type Repository interface {
	RegisterProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProductsBy(ctx context.Context, p *Product) ([]Product, error)
	GetAllProductsByCategory(ctx context.Context, category string) ([]Product, error)
}

type RepositoryCategory interface {
	RegisterCategoryProduct(ctx context.Context, category *CategoryProduct) error
	UpdateCategoryProduct(ctx context.Context, category *CategoryProduct) error
	DeleteCategoryProduct(ctx context.Context, id string) error
	GetCategoryProductById(ctx context.Context, id string) (*CategoryProduct, error)
	GetAllCategoryProduct(ctx context.Context) ([]CategoryProduct, error)
}
