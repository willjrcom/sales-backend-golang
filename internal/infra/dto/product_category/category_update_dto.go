package productcategorydto

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

var ()

type CategoryUpdateDTO struct {
	Name                 *string               `json:"name"`
	ImagePath            *string               `json:"image_path"`
	NeedPrint            *bool                 `json:"need_print"`
	PrinterName          string                `json:"printer_name"`
	UseProcessRule       *bool                 `json:"use_process_rule"`
	RemovableIngredients []string              `json:"removable_ingredients"`
	AdditionalCategories []entitydto.IDRequest `json:"additional_categories"`
	ComplementCategories []entitydto.IDRequest `json:"complement_categories"`
}

func (c *CategoryUpdateDTO) UpdateDomain(category *productentity.ProductCategory) (err error) {
	if c.Name != nil {
		category.Name = *c.Name
	}

	if c.ImagePath != nil {
		category.ImagePath = *c.ImagePath
	}

	if c.NeedPrint != nil {
		category.NeedPrint = *c.NeedPrint
		category.PrinterName = c.PrinterName
	}

	if c.UseProcessRule != nil {
		category.UseProcessRule = *c.UseProcessRule
	}

	if len(c.RemovableIngredients) != 0 {
		category.RemovableIngredients = c.RemovableIngredients
	}

	if c.AdditionalCategories != nil {
		category.AdditionalCategories = make([]productentity.ProductCategory, len(c.AdditionalCategories))
		for i := range c.AdditionalCategories {
			category.AdditionalCategories[i] = productentity.ProductCategory{
				Entity: entity.Entity{
					ID: c.AdditionalCategories[i].ID,
				},
			}
		}
	}

	if c.ComplementCategories != nil {
		category.ComplementCategories = make([]productentity.ProductCategory, len(c.ComplementCategories))
		for i := range c.ComplementCategories {
			category.ComplementCategories[i] = productentity.ProductCategory{
				Entity: entity.Entity{
					ID: c.ComplementCategories[i].ID,
				},
			}
		}
	}

	return nil
}
