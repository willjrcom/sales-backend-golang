package productdto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameAndActiveIsEmpty = errors.New("name and active can't be empty")
	ErrInvalidInput         = errors.New("invalid input")
)

type UpdateSizeInput struct {
	Name   *string `json:"name"`
	Active *bool   `json:"active"`
}

func (s *UpdateSizeInput) validate() error {
	if s.Name == nil && s.Active == nil {
		return ErrNameAndActiveIsEmpty
	}
	return nil
}
func (s *UpdateSizeInput) UpdateModel(model *productentity.Size) error {
	if s.validate() != nil {
		return ErrInvalidInput
	}

	if s.Name != nil {
		model.Name = *s.Name
	}
	if s.Active != nil {
		model.Active = *s.Active
	}

	return nil
}
