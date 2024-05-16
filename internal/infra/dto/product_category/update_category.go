package productcategorydto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var ()

type UpdateCategoryInput struct {
	productentity.PatchProductCategory
}

func (c *UpdateCategoryInput) UpdateModel(category *productentity.ProductCategory) (err error) {
	if c.Name != nil {
		category.Name = *c.Name
	}

	if c.ImagePath != nil {
		category.ImagePath = *c.ImagePath
	}

	if len(c.RemovableIngredients) != 0 {
		category.RemovableIngredients = c.RemovableIngredients
	}

	if c.AdditionalCategories != nil {
		category.AdditionalCategories = c.AdditionalCategories
	}

	if c.ComplementCategories != nil {
		category.ComplementCategories = c.ComplementCategories
	}

	return nil
}
