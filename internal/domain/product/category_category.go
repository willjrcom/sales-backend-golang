package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductCategoryToAdditional struct {
	bun.BaseModel        `bun:"table:product_category_to_additional"`
	CategoryID           uuid.UUID        `bun:"type:uuid,pk"`
	Category             *ProductCategory `bun:"rel:belongs-to,join:category_id=id"`
	AdditionalCategoryID uuid.UUID        `bun:"type:uuid,pk"`
	AdditionalCategory   *ProductCategory `bun:"rel:belongs-to,join:additional_category_id=id"`
}

type CategoryRelation struct {
	ID uuid.UUID `json:"id"`
}
