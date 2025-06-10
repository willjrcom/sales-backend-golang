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
	PrinterName          string
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

func NewProductCategory(categoryCommonAttributes ProductCategoryCommonAttributes) *ProductCategory {
	return &ProductCategory{
		Entity:                          entity.NewEntity(),
		ProductCategoryCommonAttributes: categoryCommonAttributes,
	}
}
