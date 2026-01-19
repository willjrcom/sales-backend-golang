package quantityrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errQuantityExists   = errors.New("quantity already exists")
	errQuantityNotFound = errors.New("quantity not found")
)

type QuantityRepositoryLocal struct {
	quantitys map[uuid.UUID]*model.Quantity
}

func NewQuantityRepositoryLocal() model.QuantityRepository {
	return &QuantityRepositoryLocal{quantitys: make(map[uuid.UUID]*model.Quantity)}
}

func (r *QuantityRepositoryLocal) CreateQuantity(_ context.Context, p *model.Quantity) error {

	if _, ok := r.quantitys[p.ID]; ok {
		return errQuantityExists
	}

	r.quantitys[p.ID] = p
	return nil
}

func (r *QuantityRepositoryLocal) UpdateQuantity(_ context.Context, s *model.Quantity) error {
	r.quantitys[s.ID] = s
	return nil
}

func (r *QuantityRepositoryLocal) DeleteQuantity(_ context.Context, id string) error {

	if _, ok := r.quantitys[uuid.MustParse(id)]; !ok {
		return errQuantityNotFound
	}

	delete(r.quantitys, uuid.MustParse(id))
	return nil
}

func (r *QuantityRepositoryLocal) GetQuantityById(_ context.Context, id string) (*model.Quantity, error) {

	if p, ok := r.quantitys[uuid.MustParse(id)]; ok {
		return p, nil
	}

	return nil, errQuantityNotFound
}

func (r *QuantityRepositoryLocal) GetQuantitiesByCategoryId(_ context.Context, categoryId string) ([]*model.Quantity, error) {
	quantities := []*model.Quantity{}
	for _, quantity := range r.quantitys {
		if quantity.CategoryID == uuid.MustParse(categoryId) {
			quantities = append(quantities, quantity)
		}
	}
	return quantities, nil
}
