package productcategorydto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
)

type RegisterCategoryInput struct {
	productentity.ProductCategoryCommonAttributes
}

func (c *RegisterCategoryInput) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *RegisterCategoryInput) ToModel() (*productentity.ProductCategory, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	categoryCommonAttributes := productentity.ProductCategoryCommonAttributes{
		Name:                 c.Name,
		AdditionalCategories: c.AdditionalCategories,
		RemovableIngredients: c.RemovableIngredients,
		ImagePath:            c.ImagePath,
		NeedPrint:            c.NeedPrint,
	}

	return productentity.NewProductCategory(categoryCommonAttributes), nil
}
