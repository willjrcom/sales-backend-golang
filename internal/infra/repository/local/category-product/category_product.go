package categoryrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	errCategoryProductAlreadyExists = errors.New("category Product already exists")
	errCategoryProductNotFound      = errors.New("category Product not found")
)

type CategoryRepositoryLocal struct {
	Categoryproducts map[uuid.UUID]*productentity.Category
}

func NewCategoryRepositoryLocal() *CategoryRepositoryLocal {
	return &CategoryRepositoryLocal{Categoryproducts: make(map[uuid.UUID]*productentity.Category)}
}

func (r *CategoryRepositoryLocal) RegisterCategory(_ context.Context, p *productentity.Category) error {
	if _, ok := r.Categoryproducts[p.ID]; ok {

		return errCategoryProductAlreadyExists
	}

	r.Categoryproducts[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) UpdateCategory(_ context.Context, p *productentity.Category) error {
	r.Categoryproducts[p.ID] = p
	return nil
}

func (r *CategoryRepositoryLocal) DeleteCategory(_ context.Context, id string) error {
	if _, ok := r.Categoryproducts[uuid.MustParse(id)]; !ok {

		return errCategoryProductNotFound
	}

	delete(r.Categoryproducts, uuid.MustParse(id))
	return nil
}

func (r *CategoryRepositoryLocal) GetCategoryById(_ context.Context, id string) (*productentity.Category, error) {
	if p, ok := r.Categoryproducts[uuid.MustParse(id)]; ok {

		return p, nil
	}

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
