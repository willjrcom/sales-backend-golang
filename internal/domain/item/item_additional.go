package itementity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ItemToAdditional struct {
	bun.BaseModel    `bun:"table:item_to_additional"`
	ItemID           uuid.UUID `bun:"type:uuid,pk" json:"item_id"`
	Item             *Item     `bun:"rel:belongs-to,join:item_id=id" json:"item,omitempty"`
	AdditionalItemID uuid.UUID `bun:"type:uuid,pk" json:"additional_item_id"`
	AdditionalItem   *Item     `bun:"rel:belongs-to,join:additional_item_id=id" json:"additional_item,omitempty"`
}

type ItemRelation struct {
	ID uuid.UUID `json:"id"`
}
