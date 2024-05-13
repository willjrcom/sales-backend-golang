package productcategorysizedto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameAndActiveIsEmpty = errors.New("name and active can't be empty")
)

type UpdateSizeInput struct {
	productentity.PatchSize
}

func (s *UpdateSizeInput) validate() error {
	if s.Name == nil && s.Active == nil {
		return ErrNameAndActiveIsEmpty
	}

	return nil
}
func (s *UpdateSizeInput) UpdateModel(model *productentity.Size) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Name != nil {
		model.Name = *s.Name
	}
	if s.Active != nil {
		model.Active = s.Active
	}

	return nil
}
