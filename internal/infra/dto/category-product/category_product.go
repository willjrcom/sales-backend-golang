package categoryproductdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryProduct struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Sizes []string  `json:"sizes"`
}

func (c *CategoryProduct) Validate() error {
	if c.Name == "" {
		return ErrCategoryNameRequired
	}
	if len(c.Sizes) == 0 {
		return ErrCategorySizeRequired
	}

	return nil
}

func (c *CategoryProduct) ToModel() (*productentity.CategoryProduct, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &productentity.CategoryProduct{
		Name:  c.Name,
		Sizes: c.Sizes,
	}, nil
}
