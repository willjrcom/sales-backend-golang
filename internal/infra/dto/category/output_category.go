package categorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.CategoryCommonAttributes
}

func (c *CategoryOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID

	if len(model.Sizes) == 0 {
		model.Sizes = []productentity.Size{}
	}

	if len(model.Quantities) == 0 {
		model.Quantities = []productentity.Quantity{}
	}

	if len(model.Products) == 0 {
		model.Products = []productentity.Product{}
	}

	if len(model.ProcessRules) == 0 {
		model.ProcessRules = []productentity.ProcessRule{}
	}

	if len(model.AdditionalCategories) == 0 {
		model.AdditionalCategories = []productentity.Category{}
	}

	c.CategoryCommonAttributes = model.CategoryCommonAttributes
}
