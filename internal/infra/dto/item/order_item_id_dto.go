package itemdto

import "github.com/google/uuid"

type ItemIDDTO struct {
	ItemID      uuid.UUID `json:"item_id"`
	GroupItemID uuid.UUID `json:"group_item_id"`
}

func FromDomain(itemID uuid.UUID, groupItemID uuid.UUID) *ItemIDDTO {
	return &ItemIDDTO{
		GroupItemID: groupItemID,
		ItemID:      itemID,
	}
}
