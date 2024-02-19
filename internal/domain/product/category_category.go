package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CategoryToAdditional struct {
	bun.BaseModel        `bun:"table:category_to_additional"`
	CategoryID           uuid.UUID `bun:"type:uuid,pk"`
	Category             *Category `bun:"rel:belongs-to,join:category_id=id"`
	AdditionalCategoryID uuid.UUID `bun:"type:uuid,pk"`
	AdditionalCategory   *Category `bun:"rel:belongs-to,join:additional_category_id=id"`
}

type CategoryRelation struct {
	ID uuid.UUID `json:"id"`
}
