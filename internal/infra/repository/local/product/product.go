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

func (r *ProductRepositoryLocal) GetAllProducts(_ context.Context) ([]productentity.Product, error) {
	products := make([]productentity.Product, 0)

	for _, p := range r.products {
		products = append(products, *p)
	}

	return products, nil
}
