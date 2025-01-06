package itemdto

import "github.com/google/uuid"

type ItemIDAndGroupItemDTO struct {
	ItemID      uuid.UUID `json:"item_id"`
	GroupItemID uuid.UUID `json:"group_item_id"`
}

func NewOutput(itemID uuid.UUID, groupItemID uuid.UUID) *ItemIDAndGroupItemDTO {
	return &ItemIDAndGroupItemDTO{
		GroupItemID: groupItemID,
		ItemID:      itemID,
	}
}
