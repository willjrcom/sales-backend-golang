package quantitydto

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrQuantitiesRequired = errors.New("quantities is required")
)

type QuantityCreateBatchDTO struct {
	Quantities []float64 `json:"quantities"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *QuantityCreateBatchDTO) validate() error {
	if len(s.Quantities) == 0 {
		return ErrQuantitiesRequired
	}

	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *QuantityCreateBatchDTO) ToDomain() ([]float64, *uuid.UUID, error) {
	if err := s.validate(); err != nil {
		return nil, nil, err
	}

	return s.Quantities, &s.CategoryID, nil
}
