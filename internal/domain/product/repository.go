package productentity

import (
	"context"
)

type ProductRepository interface {
	RegisterProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProductsBy(ctx context.Context, p *Product) ([]Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
}

type CategoryRepository interface {
	RegisterCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryById(ctx context.Context, id string) (*Category, error)
	GetAllCategoryProducts(ctx context.Context) ([]Category, error)
	GetAllCategorySizes(ctx context.Context) ([]Category, error)
}

type SizeRepository interface {
	RegisterSize(ctx context.Context, Size *Size) error
	UpdateSize(ctx context.Context, Size *Size) error
	DeleteSize(ctx context.Context, id string) error
	GetSizeById(ctx context.Context, id string) (*Size, error)
}

type QuantityRepository interface {
	RegisterQuantity(ctx context.Context, Quantity *Quantity) error
	UpdateQuantity(ctx context.Context, Quantity *Quantity) error
	DeleteQuantity(ctx context.Context, id string) error
	GetQuantityById(ctx context.Context, id string) (*Quantity, error)
}

type ProcessRepository interface {
	RegisterProcess(ctx context.Context, Process *Process) error
	UpdateProcess(ctx context.Context, Process *Process) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*Process, error)
}
