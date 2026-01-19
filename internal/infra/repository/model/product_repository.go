package model

import (
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProductByCode(ctx context.Context, code string) (*Product, error)
	GetAllProducts(ctx context.Context, page, perPage int, isActive bool, categoryID string) ([]Product, int, error)
	GetDefaultProducts(ctx context.Context, page, perPage int, isActive bool) ([]Product, int, error)
	GetAllProductsMap(ctx context.Context, isActive bool) ([]Product, error)
}
