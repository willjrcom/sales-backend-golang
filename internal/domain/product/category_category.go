package productentity

import "github.com/google/uuid"

type CategoryCategory struct {
	CategoryID           uuid.UUID `bun:",pk"`
	Category             *Category `bun:"rel:belongs-to,join:category_id=id"`
	CategoryAdditionalID uuid.UUID `bun:",pk"`
	CategoryAdditional   *Category `bun:"rel:belongs-to,join:category_id=id"`
}

type CategoryRelation struct {
	ID uuid.UUID `json:"id"`
}
