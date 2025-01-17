package sizedto

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSizesRequired = errors.New("sizes is required")
)

type SizeCreateBatchDTO struct {
	Sizes      []string  `json:"sizes"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *SizeCreateBatchDTO) validate() error {
	if len(s.Sizes) == 0 {
		return ErrSizesRequired
	}

	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *SizeCreateBatchDTO) ToDomain() ([]string, *uuid.UUID, error) {
	if err := s.validate(); err != nil {
		return nil, nil, err
	}

	return s.Sizes, &s.CategoryID, nil
}
