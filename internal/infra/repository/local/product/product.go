package productrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	errProductExists   = errors.New("product already exists")
	errProductNotFound = errors.New("product not found")
)

type ProductRepositoryLocal struct {
	mu       sync.Mutex
	products map[uuid.UUID]*productentity.Product
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{products: make(map[uuid.UUID]*productentity.Product)}
}

func (r *ProductRepositoryLocal) RegisterProduct(_ context.Context, p *productentity.Product) error {
	r.mu.Lock()

	if _, ok := r.products[p.ID]; ok {
		r.mu.Unlock()
		return errProductExists
	}

	r.products[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) UpdateProduct(_ context.Context, p *productentity.Product) error {
	r.mu.Lock()
	r.products[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) DeleteProduct(_ context.Context, id string) error {
	r.mu.Lock()

	if _, ok := r.products[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errProductNotFound
	}

	delete(r.products, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) GetProductById(_ context.Context, id string) (*productentity.Product, error) {
	r.mu.Lock()

	if p, ok := r.products[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errProductNotFound
}

func (r *ProductRepositoryLocal) GetProductsBy(_ context.Context, p *productentity.Product) ([]productentity.Product, error) {
	products := make([]productentity.Product, 0)

	for _, v := range r.products {
		if v.Code != "" && v.Code == p.Code {
			products = append(products, *v)
		}
		if v.Name != "" && v.Name == p.Name {
			products = append(products, *v)
		}
		if v.CategoryID != uuid.Nil && v.CategoryID == p.CategoryID {
			products = append(products, *v)
		}
		if v.Size != "" && v.Size == p.Size {
			products = append(products, *v)
		}
	}

	return products, nil
}

func (r *ProductRepositoryLocal) GetAllProductsByCategory(_ context.Context, category string) ([]productentity.Product, error) {
	products := make([]productentity.Product, 0)

	for _, p := range r.products {
		if category == "teste" {
			products = append(products, *p)
		}
	}

	return products, nil
}
