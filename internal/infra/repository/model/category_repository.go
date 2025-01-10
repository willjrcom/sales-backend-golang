package model

import "context"

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *ProductCategory) error
	UpdateCategory(ctx context.Context, category *ProductCategory) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryById(ctx context.Context, id string) (*ProductCategory, error)
	GetCategoryByName(ctx context.Context, name string, withRelation bool) (*ProductCategory, error)
	GetAllCategories(ctx context.Context) ([]ProductCategory, error)
	GetAllCategoriesWithProcessRulesAndOrderProcess(ctx context.Context) ([]ProductCategoryWithOrderProcess, error)
}
