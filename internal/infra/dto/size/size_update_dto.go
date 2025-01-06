package sizedto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameAndActiveIsEmpty = errors.New("name and active can't be empty")
)

type SizeUpdateDTO struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

func (s *SizeUpdateDTO) validate() error {
	if s.Name == nil && s.IsActive == nil {
		return ErrNameAndActiveIsEmpty
	}

	return nil
}
func (s *SizeUpdateDTO) UpdateDomain(size *productentity.Size) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Name != nil {
		size.Name = *s.Name
	}
	if s.IsActive != nil {
		size.IsActive = s.IsActive
	}

	return nil
}
