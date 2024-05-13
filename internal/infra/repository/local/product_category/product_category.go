package productcategoryrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	errProductCategoryAlreadyExists = errors.New("product category already exists")
	errProductCategoryNotFound      = errors.New("product category not found")
)

type CategoryRepositoryLocal struct {
	productCategories map[uuid.UUID]*productentity.ProductCategory
}

func NewCategoryRepositoryLocal() *CategoryRepositoryLocal {
	return &CategoryRepositoryLocal{productCategories: make(map[uuid.UUID]*productentity.ProductCategory)}
}

func (r *CategoryRepositoryLocal) CreateCategory(_ context.Context, p *productentity.ProductCategory) error {
	if _, ok := r.productCategories[p.ID]; ok {

		return errProductCategoryAlreadyExists
	}

	r.productCategories[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) UpdateCategory(_ context.Context, p *productentity.ProductCategory) error {
	r.productCategories[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) DeleteCategory(_ context.Context, id string) error {
	if _, ok := r.productCategories[uuid.MustParse(id)]; !ok {

		return errProductCategoryNotFound
	}

	delete(r.productCategories, uuid.MustParse(id))
	return nil
}

func (r *CategoryRepositoryLocal) GetCategoryById(_ context.Context, id string) (*productentity.ProductCategory, error) {
	if p, ok := r.productCategories[uuid.MustParse(id)]; ok {

		return p, nil
	}

	return nil, errProductCategoryNotFound
}

func (r *CategoryRepositoryLocal) GetCategoryByName(_ context.Context, name string, withRelation bool) (*productentity.ProductCategory, error) {
	for _, p := range r.productCategories {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, errProductCategoryNotFound
}

func (r *CategoryRepositoryLocal) GetAllCategories(_ context.Context) ([]productentity.ProductCategory, error) {
	productCategories := make([]productentity.ProductCategory, 0)

	for _, p := range r.productCategories {
		p.Sizes = nil
		productCategories = append(productCategories, *p)
	}

	return productCategories, nil
}
