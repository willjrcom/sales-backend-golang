package productcategoryrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errProductCategoryAlreadyExists = errors.New("product category already exists")
	errProductCategoryNotFound      = errors.New("product category not found")
)

type CategoryRepositoryLocal struct {
	mu                sync.RWMutex
	productCategories map[uuid.UUID]*model.ProductCategory
}

func NewCategoryRepositoryLocal() model.CategoryRepository {
	return &CategoryRepositoryLocal{productCategories: make(map[uuid.UUID]*model.ProductCategory)}
}

func (r *CategoryRepositoryLocal) CreateCategory(_ context.Context, p *model.ProductCategory) error {
	if _, ok := r.productCategories[p.ID]; ok {

		return errProductCategoryAlreadyExists
	}

	r.productCategories[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) UpdateCategory(_ context.Context, p *model.ProductCategory) error {
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

func (r *CategoryRepositoryLocal) GetCategoryById(_ context.Context, id string) (*model.ProductCategory, error) {
	if p, ok := r.productCategories[uuid.MustParse(id)]; ok {

		return p, nil
	}

	return nil, errProductCategoryNotFound
}

func (r *CategoryRepositoryLocal) GetCategoryByName(_ context.Context, name string, withRelation bool) (*model.ProductCategory, error) {
	for _, p := range r.productCategories {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, errProductCategoryNotFound
}

func (r *CategoryRepositoryLocal) GetAllCategories(_ context.Context, IDs []uuid.UUID, page int, perPage int, isActive ...bool) ([]model.ProductCategory, error) {
	productCategories := make([]model.ProductCategory, 0)

	for _, p := range r.productCategories {
		p.Sizes = nil
		productCategories = append(productCategories, *p)
	}

	return productCategories, nil
}

func (r *CategoryRepositoryLocal) GetAllCategoriesMap(_ context.Context, isActive bool, isAdditional, isComplement *bool) ([]model.ProductCategory, error) {
	productCategories := make([]model.ProductCategory, 0)

	for _, p := range r.productCategories {
		if p.IsActive != isActive {
			continue
		}

		if isAdditional != nil && p.IsAdditional != *isAdditional {
			continue
		}

		if isComplement != nil && p.IsComplement != *isComplement {
			continue
		}

		cat := model.ProductCategory{}
		cat.ID = p.ID
		cat.Name = p.Name
		productCategories = append(productCategories, cat)
	}

	return productCategories, nil
}

func (r *CategoryRepositoryLocal) GetAllCategoriesWithProcessRulesAndOrderProcess(_ context.Context) ([]model.ProductCategoryWithOrderProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.ProductCategoryWithOrderProcess, 0, len(r.productCategories))
	for _, pc := range r.productCategories {
		// Convert DB model to domain, then to combined DTO
		dom := pc.ToDomain()
		var cp model.ProductCategoryWithOrderProcess
		cp.FromDomain(dom)
		out = append(out, cp)
	}
	return out, nil
}

func (r *CategoryRepositoryLocal) GetComplementProducts(_ context.Context, categoryID string) ([]model.Product, error) {
	// Dummy implementation
	return []model.Product{}, nil
}

func (r *CategoryRepositoryLocal) GetAdditionalProducts(_ context.Context, categoryID string) ([]model.Product, error) {
	// Dummy implementation
	return []model.Product{}, nil
}

func (r *CategoryRepositoryLocal) GetComplementCategories(_ context.Context) ([]model.ProductCategory, error) {
	// Dummy implementation
	return []model.ProductCategory{}, nil
}

func (r *CategoryRepositoryLocal) GetAdditionalCategories(_ context.Context) ([]model.ProductCategory, error) {
	// Dummy implementation
	return []model.ProductCategory{}, nil
}

func (r *CategoryRepositoryLocal) GetDefaultCategories(_ context.Context) ([]model.ProductCategory, error) {
	// Dummy implementation
	return []model.ProductCategory{}, nil
}
