package productcategorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"

	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
)

type CategoryDTO struct {
	ID                   uuid.UUID                       `json:"id"`
	Name                 string                          `json:"name"`
	ImagePath            string                          `json:"image_path,omitempty"`
	NeedPrint            bool                            `json:"need_print"`
	UseProcessRule       bool                            `json:"use_process_rule"`
	IsAdditional         bool                            `json:"is_additional"`
	IsComplement         bool                            `json:"is_complement"`
	RemovableIngredients []string                        `json:"removable_ingredients,omitempty"`
	Sizes                []sizedto.SizeDTO               `json:"sizes,omitempty"`
	Quantities           []quantitydto.QuantityDTO       `json:"quantities,omitempty"`
	Products             []ProductDTO                    `json:"products,omitempty"`
	AdditionalCategories []CategoryDTO                   `json:"additional_categories,omitempty"`
	ComplementCategories []CategoryDTO                   `json:"complement_categories,omitempty"`
	ProcessRules         []processruledto.ProcessRuleDTO `json:"process_rules,omitempty"`
}

func (c *CategoryDTO) FromDomain(category *productentity.ProductCategory) {
	if category == nil {
		return
	}
	*c = CategoryDTO{
		ID:                   category.ID,
		Name:                 category.Name,
		ImagePath:            category.ImagePath,
		NeedPrint:            category.NeedPrint,
		UseProcessRule:       category.UseProcessRule,
		IsAdditional:         category.IsAdditional,
		IsComplement:         category.IsComplement,
		RemovableIngredients: category.RemovableIngredients,
	}

	if len(category.Sizes) == 0 {
		category.Sizes = []productentity.Size{}
	}

	if len(category.Quantities) == 0 {
		category.Quantities = []productentity.Quantity{}
	}

	if len(category.Products) == 0 {
		category.Products = []productentity.Product{}
	}

	if len(category.ProcessRules) == 0 {
		category.ProcessRules = []productentity.ProcessRule{}
	} else {
		c.ProcessRules = make([]processruledto.ProcessRuleDTO, len(category.ProcessRules))

		for i, processRule := range category.ProcessRules {
			c.ProcessRules[i] = processruledto.ProcessRuleDTO{}
			c.ProcessRules[i].FromDomain(&processRule)
		}
	}

	if len(category.AdditionalCategories) == 0 {
		category.AdditionalCategories = []productentity.ProductCategory{}
	}

	if len(category.ComplementCategories) == 0 {
		category.ComplementCategories = []productentity.ProductCategory{}
	}

	for i := range category.Sizes {
		c.Sizes[i] = sizedto.SizeDTO{}
		c.Sizes[i].FromDomain(&category.Sizes[i])
	}

	for i := range category.Quantities {
		c.Quantities[i] = quantitydto.QuantityDTO{}
		c.Quantities[i].FromDomain(&category.Quantities[i])
	}

	for i := range category.Products {
		c.Products[i] = ProductDTO{}
		c.Products[i] = *FromDomain(&category.Products[i])
	}

	for i := range category.AdditionalCategories {
		c.AdditionalCategories[i] = CategoryDTO{}
		c.AdditionalCategories[i].FromDomain(&category.AdditionalCategories[i])
	}

	for i := range category.ComplementCategories {
		c.ComplementCategories[i] = CategoryDTO{}
		c.ComplementCategories[i].FromDomain(&category.ComplementCategories[i])
	}
}
