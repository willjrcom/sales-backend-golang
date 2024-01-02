package itementity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Item struct {
	entity.Entity
	bun.BaseModel `bun:"table:items"`
	ItemCommonAttributes
}

type ItemCommonAttributes struct {
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	Observation string    `bun:"observation" json:"observation"`
	Price       float64   `bun:"price,notnull" json:"price"`
	Quantity    float64   `bun:"quantity,notnull" json:"quantity"`
	GroupItemID uuid.UUID `bun:"group_item_id,type:uuid,notnull" json:"group_item_id"`
}
