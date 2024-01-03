package categoryrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	errCategoryProductAlreadyExists = errors.New("category Product already exists")
	errCategoryProductNotFound      = errors.New("category Product not found")
)

type CategoryRepositoryLocal struct {
	mu               sync.Mutex
	Categoryproducts map[uuid.UUID]*productentity.Category
}

func NewCategoryRepositoryLocal() *CategoryRepositoryLocal {
	return &CategoryRepositoryLocal{Categoryproducts: make(map[uuid.UUID]*productentity.Category)}
}

func (r *CategoryRepositoryLocal) RegisterCategory(_ context.Context, p *productentity.Category) error {
	r.mu.Lock()

	if _, ok := r.Categoryproducts[p.ID]; ok {
		r.mu.Unlock()
		return errCategoryProductAlreadyExists
	}

	r.Categoryproducts[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *CategoryRepositoryLocal) UpdateCategory(_ context.Context, p *productentity.Category) error {
	r.mu.Lock()
	r.Categoryproducts[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *CategoryRepositoryLocal) DeleteCategory(_ context.Context, id string) error {
	r.mu.Lock()

	if _, ok := r.Categoryproducts[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errCategoryProductNotFound
	}

	delete(r.Categoryproducts, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *CategoryRepositoryLocal) GetCategoryById(_ context.Context, id string) (*productentity.Category, error) {
	r.mu.Lock()

	if p, ok := r.Categoryproducts[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errCategoryProductNotFound
}

func (r *CategoryRepositoryLocal) GetAllCategories(_ context.Context) ([]productentity.Category, error) {
	Categoryproducts := make([]productentity.Category, 0)

	for _, p := range r.Categoryproducts {
		p.Sizes = nil
		Categoryproducts = append(Categoryproducts, *p)
	}

	return Categoryproducts, nil
}
