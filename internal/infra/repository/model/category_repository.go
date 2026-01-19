package model

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *ProductCategory) error
	UpdateCategory(ctx context.Context, category *ProductCategory) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryById(ctx context.Context, id string) (*ProductCategory, error)
	GetCategoryByName(ctx context.Context, name string, withRelation bool) (*ProductCategory, error)
	GetAllCategories(ctx context.Context, IDs []uuid.UUID, page int, perPage int, isActive ...bool) ([]ProductCategory, error)
	GetAllCategoriesMap(ctx context.Context, isActive bool) ([]ProductCategory, error)
	GetAllCategoriesWithProcessRulesAndOrderProcess(ctx context.Context) ([]ProductCategoryWithOrderProcess, error)
	GetComplementProducts(ctx context.Context, categoryID string) ([]Product, error)
	GetAdditionalProducts(ctx context.Context, categoryID string) ([]Product, error)
	GetComplementCategories(ctx context.Context) ([]ProductCategory, error)
	GetAdditionalCategories(ctx context.Context) ([]ProductCategory, error)
	GetDefaultCategories(ctx context.Context) ([]ProductCategory, error)
}
