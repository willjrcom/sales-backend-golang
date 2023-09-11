package productdto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
	ErrSizeIsEmpty = errors.New("sizes is empty")
)

type RegisterCategoryProductInput struct {
	Name  string   `json:"name"`
	Sizes []string `json:"sizes"`
}

func (c *RegisterCategoryProductInput) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}
	if len(c.Sizes) == 0 {
		return ErrSizeIsEmpty
	}

	return nil
}

func (c *RegisterCategoryProductInput) ToModel() (*productentity.CategoryProduct, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return &productentity.CategoryProduct{
		Entity: entity.NewEntity(),
		Name:   c.Name,
		Sizes:  c.Sizes,
	}, nil
}
