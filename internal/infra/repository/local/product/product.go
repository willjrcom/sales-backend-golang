package productrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errProductExists   = errors.New("product already exists")
	errProductNotFound = errors.New("product not found")
)

type ProductRepositoryLocal struct {
	mu       sync.Mutex
	products map[uuid.UUID]*model.Product
}

func NewProductRepositoryLocal() model.ProductRepository {
	return &ProductRepositoryLocal{products: make(map[uuid.UUID]*model.Product)}
}

func (r *ProductRepositoryLocal) CreateProduct(_ context.Context, p *model.Product) error {

	if _, ok := r.products[p.ID]; ok {

		return errProductExists
	}

	r.products[p.ID] = p
	return nil
}

func (r *ProductRepositoryLocal) UpdateProduct(_ context.Context, p *model.Product) error {
	r.products[p.ID] = p
	return nil
}

func (r *ProductRepositoryLocal) DeleteProduct(_ context.Context, id string) error {

	if _, ok := r.products[uuid.MustParse(id)]; !ok {

		return errProductNotFound
	}

	delete(r.products, uuid.MustParse(id))
	return nil
}

func (r *ProductRepositoryLocal) GetProductById(_ context.Context, id string) (*model.Product, error) {

	if p, ok := r.products[uuid.MustParse(id)]; ok {

		return p, nil
	}

	return nil, errProductNotFound
}

func (r *ProductRepositoryLocal) GetProductBySKU(_ context.Context, sku string) (*model.Product, error) {

	for _, p := range r.products {
		if p.SKU == sku {
			return p, nil
		}
	}

	return nil, errProductNotFound
}

func (r *ProductRepositoryLocal) GetAllProducts(_ context.Context, page, perPage int, isActive bool, categoryID string) ([]model.Product, int, error) {
	products := make([]model.Product, 0)

	// Filter by isActive
	for _, p := range r.products {
		if p.IsActive == isActive {
			products = append(products, *p)
		}
	}

	total := len(products)

	// Apply pagination
	start := page * perPage
	end := start + perPage
	if start > total {
		return []model.Product{}, total, nil
	}
	if end > total {
		end = total
	}

	return products[start:end], total, nil
}

func (r *ProductRepositoryLocal) GetDefaultProducts(_ context.Context, page, perPage int, isActive bool) ([]model.Product, int, error) {
	products := make([]model.Product, 0)

	// Filter by isActive
	for _, p := range r.products {
		if p.IsActive == isActive {
			products = append(products, *p)
		}
	}

	total := len(products)

	// Apply pagination
	start := page * perPage
	end := start + perPage
	if start > total {
		return []model.Product{}, total, nil
	}
	if end > total {
		end = total
	}

	return products[start:end], total, nil
}

func (r *ProductRepositoryLocal) GetAllProductsMap(_ context.Context, isActive bool, categoryID string) ([]model.Product, error) {
	products := make([]model.Product, 0)

	for _, p := range r.products {
		if p.IsActive == isActive {
			products = append(products, *p)
		}
	}

	return products, nil
}
