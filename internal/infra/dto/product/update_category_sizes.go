package productdto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrSizesIsEmpty = errors.New("sizes are empty")
)

type UpdateCategoryProductSizesInput struct {
	Sizes []string `json:"sizes"`
}

func (c *UpdateCategoryProductSizesInput) validate() error {
	if len(c.Sizes) == 0 {
		return ErrSizesIsEmpty
	}

	return nil
}

func (c *UpdateCategoryProductSizesInput) UpdateModel(category *productentity.CategoryProduct) error {
	if err := c.validate(); err != nil {
		return err
	}
	if len(c.Sizes) > 0 {
		category.Sizes = c.Sizes
	}

	return nil
}
