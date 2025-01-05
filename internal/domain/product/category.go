package productentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProductCategory struct {
	entity.Entity
	ProductCategoryCommonAttributes
}

type ProductCategoryCommonAttributes struct {
	Name                 string
	ImagePath            string
	NeedPrint            bool
	UseProcessRule       bool
	RemovableIngredients []string
	Sizes                []Size
	Quantities           []Quantity
	Products             []Product
	ProcessRules         []ProcessRule
	IsAdditional         bool
	IsComplement         bool
	AdditionalCategories []ProductCategory
	ComplementCategories []ProductCategory
}

type PatchProductCategory struct {
	Name                 *string           `json:"name"`
	ImagePath            *string           `json:"image_path"`
	NeedPrint            *bool             `json:"need_print"`
	UseProcessRule       *bool             `json:"use_process_rule"`
	RemovableIngredients []string          `json:"removable_ingredients"`
	AdditionalCategories []ProductCategory `json:"additional_categories"`
	ComplementCategories []ProductCategory `json:"complement_categories"`
}

func NewProductCategory(categoryCommonAttributes ProductCategoryCommonAttributes) *ProductCategory {
	return &ProductCategory{
		Entity:                          entity.NewEntity(),
		ProductCategoryCommonAttributes: categoryCommonAttributes,
	}
}
