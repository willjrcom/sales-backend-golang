package productcategorysizedto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrCategoryRequired = errors.New("category is required")
)

type CreateSizeInput struct {
	Name       string    `json:"name"`
	Active     *bool     `json:"active"`
	CategoryID uuid.UUID `json:"category_id"`
}

func (s *CreateSizeInput) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}
	if s.Active == nil {
		s.Active = new(bool)
		*s.Active = true
	}

	return nil
}

func (s *CreateSizeInput) ToModel() (*productentity.Size, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	sizeCommonAttributes := productentity.SizeCommonAttributes{
		Name:       s.Name,
		Active:     s.Active,
		CategoryID: s.CategoryID,
	}

	return &productentity.Size{
		Entity:               entity.NewEntity(),
		SizeCommonAttributes: sizeCommonAttributes,
	}, nil
}
