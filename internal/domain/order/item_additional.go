package orderentity

import (
	"github.com/google/uuid"
)

type ItemToAdditional struct {
	ItemID           uuid.UUID
	Item             *Item
	AdditionalItemID uuid.UUID
	AdditionalItem   *Item
	ProductID        uuid.UUID
}

type ItemRelation struct {
	ID uuid.UUID `json:"id"`
}
