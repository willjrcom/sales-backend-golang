package productrepositorylocal

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductRepositoryLocal struct {
	mu       sync.Mutex
	products map[uuid.UUID]*productentity.Product
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{products: make(map[uuid.UUID]*productentity.Product)}
}

func (r *ProductRepositoryLocal) RegisterProduct(p *productentity.Product) error {
	r.mu.Lock()

	if _, ok := r.products[p.ID]; ok {
		r.mu.Unlock()
		return errors.New("Product already exists")
	}

	r.products[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) UpdateProduct(p *productentity.Product) error {
	r.mu.Lock()
	r.products[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) DeleteProduct(id string) error {
	r.mu.Lock()

	if _, ok := r.products[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errors.New("Product not found")
	}

	delete(r.products, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *ProductRepositoryLocal) GetProductById(id string) (*productentity.Product, error) {
	r.mu.Lock()

	if p, ok := r.products[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errors.New("Product not found")
}

func (r *ProductRepositoryLocal) GetProductBy(key string, value string) (*productentity.Product, error) {
	for _, p := range r.products {
		if key == "name" && p.Name == value {
			return p, nil
		}
		if key == "code" && p.Code == value {
			return p, nil
		}
		if key == "category" && p.Category.Name == value {
			return p, nil
		}
		if key == "size" && p.Size == value {
			return p, nil
		}
	}
	return nil, errors.New("Product not found")
}

func (r *ProductRepositoryLocal) GetAllProduct(key string, value string) ([]productentity.Product, error) {
	products := make([]productentity.Product, 0)

	for _, p := range r.products {
		products = append(products, *p)
	}

	return products, nil
}
