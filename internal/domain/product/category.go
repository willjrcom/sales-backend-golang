package productentity

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProductCategory struct {
	entity.Entity
	bun.BaseModel `bun:"table:product_categories"`
	ProductCategoryCommonAttributes
}

type ProductCategoryCommonAttributes struct {
	Name                 string            `bun:"name,unique,notnull" json:"name"`
	ImagePath            string            `bun:"image_path" json:"image_path"`
	NeedPrint            bool              `bun:"need_print,notnull" json:"need_print"`
	RemovableIngredients []string          `bun:"removable_ingredients,type:jsonb" json:"removable_ingredients,omitempty"`
	Sizes                []Size            `bun:"rel:has-many,join:id=category_id" json:"sizes,omitempty"`
	Quantities           []Quantity        `bun:"rel:has-many,join:id=category_id" json:"quantities,omitempty"`
	Products             []Product         `bun:"rel:has-many,join:id=category_id" json:"products,omitempty"`
	ProcessRules         []ProcessRule     `bun:"rel:has-many,join:id=category_id" json:"process_rules,omitempty"`
	AdditionalCategories []ProductCategory `bun:"m2m:category_to_additional,join:Category=AdditionalCategory" json:"category_to_additional,omitempty"`
}

type PatchProductCategory struct {
	Name                 *string           `json:"name"`
	ImagePath            *string           `json:"image_path"`
	NeedPrint            *bool             `json:"need_print"`
	RemovableIngredients []string          `json:"removable_ingredients"`
	AdditionalCategories []ProductCategory `json:"category_to_additional,omitempty"`
}

func NewProductCategory(categoryCommonAttributes ProductCategoryCommonAttributes) *ProductCategory {
	return &ProductCategory{
		Entity:                          entity.NewEntity(),
		ProductCategoryCommonAttributes: categoryCommonAttributes,
	}
}
