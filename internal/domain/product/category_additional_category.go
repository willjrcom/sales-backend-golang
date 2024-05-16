package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductCategoryToAdditional struct {
	bun.BaseModel        `bun:"table:product_category_to_additional"`
	CategoryID           uuid.UUID        `bun:"type:uuid,pk" json:"category_id"`
	Category             *ProductCategory `bun:"rel:belongs-to,join:category_id=id" json:"category,omitempty"`
	AdditionalCategoryID uuid.UUID        `bun:"type:uuid,pk" json:"additional_category_id"`
	AdditionalCategory   *ProductCategory `bun:"rel:belongs-to,join:additional_category_id=id" json:"additional_category,omitempty"`
}

type CategoryRelation struct {
	ID uuid.UUID `json:"id"`
}
