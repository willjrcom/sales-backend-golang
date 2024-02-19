package categorydto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
)

type RegisterCategoryInput struct {
	productentity.CategoryCommonAttributes
}

func (c *RegisterCategoryInput) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *RegisterCategoryInput) ToModel() (*productentity.Category, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	categoryCommonAttributes := productentity.CategoryCommonAttributes{
		Name:                 c.Name,
		AdditionalCategories: c.AdditionalCategories,
		RemovableItems:       c.RemovableItems,
		ImagePath:            c.ImagePath,
		NeedPrint:            c.NeedPrint,
	}

	return productentity.NewCategory(categoryCommonAttributes), nil
}
