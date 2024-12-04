package productcategorydto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
)

type CreateCategoryInput struct {
	productentity.ProductCategoryCommonAttributes
}

func (c *CreateCategoryInput) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *CreateCategoryInput) ToModel() (*productentity.ProductCategory, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return productentity.NewProductCategory(c.ProductCategoryCommonAttributes), nil
}
