package categoryrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errCategoryExists   = errors.New("category already exists")
	errCategoryNotFound = errors.New("category not found")
)

type CategoryRepositoryLocal struct {
	categories map[uuid.UUID]*model.ProductCategory
}

func NewCategoryRepositoryLocal() model.CategoryRepository {
	return &CategoryRepositoryLocal{categories: make(map[uuid.UUID]*model.ProductCategory)}
}

func (r *CategoryRepositoryLocal) CreateCategory(_ context.Context, p *model.ProductCategory) error {

	if _, ok := r.categories[p.ID]; ok {
		return errCategoryExists
	}

	r.categories[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) UpdateCategory(_ context.Context, s *model.ProductCategory) error {
	r.categories[s.ID] = s
	return nil
}

func (r *CategoryRepositoryLocal) DeleteCategory(_ context.Context, id string) error {

	if _, ok := r.categories[uuid.MustParse(id)]; !ok {
		return errCategoryNotFound
	}

	delete(r.categories, uuid.MustParse(id))
	return nil
}

func (r *CategoryRepositoryLocal) GetCategoryById(_ context.Context, id string) (*model.ProductCategory, error) {

	if p, ok := r.categories[uuid.MustParse(id)]; ok {
		return p, nil
	}

	return nil, errCategoryNotFound
}

func (r *CategoryRepositoryLocal) GetCategoryByName(_ context.Context, name string, withRelation bool) (*model.ProductCategory, error) {

	for _, p := range r.categories {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, errCategoryNotFound
}
func (r *CategoryRepositoryLocal) GetAllCategories(_ context.Context) ([]model.ProductCategory, error) {
	mapCategories := make([]model.ProductCategory, 0)

	for _, p := range r.categories {
		mapCategories = append(mapCategories, *p)
	}

	return mapCategories, nil
}

func (r *CategoryRepositoryLocal) GetAllCategoriesWithProcessRulesAndOrderProcess(_ context.Context) ([]model.ProductCategoryWithOrderProcess, error) {
	return nil, nil
}
