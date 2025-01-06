package model

import (
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProductCategory struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:product_categories"`
	ProductCategoryCommonAttributes
}

type ProductCategoryCommonAttributes struct {
	Name                 string            `bun:"name,notnull"`
	ImagePath            string            `bun:"image_path"`
	NeedPrint            bool              `bun:"need_print,notnull"`
	UseProcessRule       bool              `bun:"use_process_rule,notnull"`
	RemovableIngredients []string          `bun:"removable_ingredients,type:jsonb"`
	Sizes                []Size            `bun:"rel:has-many,join:id=category_id"`
	Quantities           []Quantity        `bun:"rel:has-many,join:id=category_id"`
	Products             []Product         `bun:"rel:has-many,join:id=category_id"`
	ProcessRules         []ProcessRule     `bun:"rel:has-many,join:id=category_id"`
	IsAdditional         bool              `bun:"is_additional"`
	IsComplement         bool              `bun:"is_complement"`
	AdditionalCategories []ProductCategory `bun:"m2m:product_category_to_additional,join:Category=AdditionalCategory"`
	ComplementCategories []ProductCategory `bun:"m2m:product_category_to_complement,join:Category=ComplementCategory"`
}
