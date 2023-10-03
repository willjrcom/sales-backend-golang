package orderdto

import (
	"github.com/google/uuid"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type AddItemOrderInput struct {
	ItemID      uuid.UUID `json:"item_id"`
	Quantity    float64   `json:"quantity"`
	Observation string    `json:"observation"`
}

func (a *AddItemOrderInput) ToModel() (uuid.UUID, *itementity.Item) {
	return a.ItemID, &itementity.Item{
		Quantity:    a.Quantity,
		Observation: a.Observation,
	}
}
