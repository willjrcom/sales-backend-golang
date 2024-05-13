package productentity

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProductByCode(ctx context.Context, code string) (*Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *ProductCategory) error
	UpdateCategory(ctx context.Context, category *ProductCategory) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryById(ctx context.Context, id string) (*ProductCategory, error)
	GetCategoryByName(ctx context.Context, name string, withRelation bool) (*ProductCategory, error)
	GetAllCategories(ctx context.Context) ([]ProductCategory, error)
}

type SizeRepository interface {
	CreateSize(ctx context.Context, Size *Size) error
	UpdateSize(ctx context.Context, Size *Size) error
	DeleteSize(ctx context.Context, id string) error
	GetSizeById(ctx context.Context, id string) (*Size, error)
}

type QuantityRepository interface {
	CreateQuantity(ctx context.Context, Quantity *Quantity) error
	UpdateQuantity(ctx context.Context, Quantity *Quantity) error
	DeleteQuantity(ctx context.Context, id string) error
	GetQuantityById(ctx context.Context, id string) (*Quantity, error)
}

type ProcessRuleRepository interface {
	CreateProcessRule(ctx context.Context, ProcessRule *ProcessRule) error
	UpdateProcessRule(ctx context.Context, ProcessRule *ProcessRule) error
	DeleteProcessRule(ctx context.Context, id string) error
	GetProcessRuleById(ctx context.Context, id string) (*ProcessRule, error)
	GetProcessRuleByCategoryIdAndOrder(ctx context.Context, id string, order int8) (*ProcessRule, error)
	GetFirstProcessRuleByCategoryId(ctx context.Context, id string) (*ProcessRule, error)
	GetMapProcessRulesByFirstOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	GetMapProcessRulesByLastOrder(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	IsLastProcessRuleByID(ctx context.Context, id uuid.UUID) (bool, error)
}
