package categoryproductdto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategorySizeRequired = errors.New("category size is required")
)

type CreateCategoryProduct struct {
	Name  string   `json:"name"`
	Sizes []string `json:"sizes"`
}

func (c *CreateCategoryProduct) validate() error {
	if c.Name == "" {
		return ErrCategoryNameRequired
	}
	if len(c.Sizes) == 0 {
		return ErrCategorySizeRequired
	}

	return nil
}

func (c *CreateCategoryProduct) ToModel() (*productentity.CategoryProduct, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return &productentity.CategoryProduct{
		Name:  c.Name,
		Sizes: c.Sizes,
	}, nil
}
