package itementity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Item struct {
	entity.Entity
	GroupItemID uuid.UUID  `bun:"group_item_id,type:uuid,notnull"`
	Name        string     `bun:"name"`
	Quantity    float64    `bun:"quantity"`
	Description string     `bun:"description"`
	Price       float64    `bun:"price"`
	Status      StatusItem `bun:"status"`
	Observation string     `bun:"observation"`
}
