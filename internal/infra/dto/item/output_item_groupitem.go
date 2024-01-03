package itemdto

import "github.com/google/uuid"

type ItemIDAndGroupItemOutput struct {
	ItemID      uuid.UUID `json:"item_id"`
	GroupItemID uuid.UUID `json:"group_item_id"`
}
