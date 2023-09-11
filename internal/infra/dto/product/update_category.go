package productdto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameAndSizesIsEmpty = errors.New("name and sizes are empty")
)

type UpdateCategoryProductInput struct {
	Name  *string  `json:"name"`
	Sizes []string `json:"sizes"`
}

func (c *UpdateCategoryProductInput) validate() error {
	if c.Name == nil && len(c.Sizes) == 0 {
		return ErrNameAndSizesIsEmpty
	}

	return nil
}

func (c *UpdateCategoryProductInput) UpdateModel(category *productentity.CategoryProduct) error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.Name != nil {
		category.Name = *c.Name
	}
	if len(c.Sizes) > 0 {
		category.Sizes = c.Sizes
	}

	return nil
}
