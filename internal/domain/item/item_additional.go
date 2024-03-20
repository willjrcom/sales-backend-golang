package itementity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ItemToAdditional struct {
	bun.BaseModel    `bun:"table:item_to_additional"`
	ItemID           uuid.UUID `bun:"type:uuid,pk"`
	Item             *Item     `bun:"rel:belongs-to,join:item_id=id"`
	AdditionalItemID uuid.UUID `bun:"type:uuid,pk"`
	AdditionalItem   *Item     `bun:"rel:belongs-to,join:additional_item_id=id"`
}

type ItemRelation struct {
	ID uuid.UUID `json:"id"`
}
