package productentity

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProductCategory struct {
	entity.Entity
	bun.BaseModel `bun:"table:product_categories"`
	ProductCategoryCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type ProductCategoryCommonAttributes struct {
	Name                 string            `bun:"name,notnull" json:"name"`
	ImagePath            string            `bun:"image_path" json:"image_path"`
	NeedPrint            bool              `bun:"need_print,notnull" json:"need_print"`
	UseProcessRule       bool              `bun:"use_process_rule,notnull" json:"use_process_rule"`
	RemovableIngredients []string          `bun:"removable_ingredients,type:jsonb" json:"removable_ingredients,omitempty"`
	Sizes                []Size            `bun:"rel:has-many,join:id=category_id" json:"sizes,omitempty"`
	Quantities           []Quantity        `bun:"rel:has-many,join:id=category_id" json:"quantities,omitempty"`
	Products             []Product         `bun:"rel:has-many,join:id=category_id" json:"products,omitempty"`
	ProcessRules         []ProcessRule     `bun:"rel:has-many,join:id=category_id" json:"process_rules,omitempty"`
	IsAdditional         bool              `bun:"is_additional" json:"is_additional"`
	IsComplement         bool              `bun:"is_complement" json:"is_complement"`
	AdditionalCategories []ProductCategory `bun:"m2m:product_category_to_additional,join:Category=AdditionalCategory" json:"additional_categories,omitempty"`
	ComplementCategories []ProductCategory `bun:"m2m:product_category_to_complement,join:Category=ComplementCategory" json:"complement_categories,omitempty"`
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
