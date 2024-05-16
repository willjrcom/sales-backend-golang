package productcategorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	productcategoryprocessruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category_process_rule"
)

type CategoryOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.ProductCategoryCommonAttributes
	ProcessRules []productcategoryprocessruledto.ProcessRuleOutput `json:"process_rules,omitempty"`
}

func (c *CategoryOutput) FromModel(model *productentity.ProductCategory) {
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
	} else {
		c.ProcessRules = make([]productcategoryprocessruledto.ProcessRuleOutput, len(model.ProcessRules))

		for i, processRule := range model.ProcessRules {
			c.ProcessRules[i] = productcategoryprocessruledto.ProcessRuleOutput{}
			c.ProcessRules[i].FromModel(&processRule)
		}
	}

	if len(model.AdditionalCategories) == 0 {
		model.AdditionalCategories = []productentity.ProductCategory{}
	}

	if len(model.ComplementCategories) == 0 {
		model.ComplementCategories = []productentity.ProductCategory{}
	}

	c.ProductCategoryCommonAttributes = model.ProductCategoryCommonAttributes
}
