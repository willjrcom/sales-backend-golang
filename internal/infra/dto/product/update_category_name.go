package productdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var ()

type UpdateCategoryProductNameInput struct {
	Name *string `json:"name"`
}

func (c *UpdateCategoryProductNameInput) validate() error {
	if c.Name == nil {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *UpdateCategoryProductNameInput) UpdateModel(category *productentity.CategoryProduct) error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.Name != nil {
		category.Name = *c.Name
	}

	return nil
}
