package categoryrepositorylocal

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryProductRepositoryLocal struct {
	mu               sync.Mutex
	Categoryproducts map[uuid.UUID]*productentity.CategoryProduct
}

func NewCategoryProductRepositoryLocal() *CategoryProductRepositoryLocal {
	return &CategoryProductRepositoryLocal{Categoryproducts: make(map[uuid.UUID]*productentity.CategoryProduct)}
}

func (r *CategoryProductRepositoryLocal) RegisterCategoryProduct(p *productentity.CategoryProduct) error {
	r.mu.Lock()

	if _, ok := r.Categoryproducts[p.ID]; ok {
		r.mu.Unlock()
		return errors.New("Category Product already exists")
	}

	r.Categoryproducts[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *CategoryProductRepositoryLocal) UpdateCategoryProduct(p *productentity.CategoryProduct) error {
	r.mu.Lock()
	r.Categoryproducts[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *CategoryProductRepositoryLocal) DeleteCategoryProduct(id string) error {
	r.mu.Lock()

	if _, ok := r.Categoryproducts[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errors.New("Category Product not found")
	}

	delete(r.Categoryproducts, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *CategoryProductRepositoryLocal) GetCategoryProductById(id string) (*productentity.CategoryProduct, error) {
	r.mu.Lock()

	if p, ok := r.Categoryproducts[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errors.New("Category Product not found")
}

func (r *CategoryProductRepositoryLocal) GetAllCategoryProduct() ([]productentity.CategoryProduct, error) {
	Categoryproducts := make([]productentity.CategoryProduct, 0)

	for _, p := range r.Categoryproducts {
		Categoryproducts = append(Categoryproducts, *p)
	}

	return Categoryproducts, nil
}
