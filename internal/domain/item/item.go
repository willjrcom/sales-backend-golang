package itementity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Item struct {
	entity.Entity
	bun.BaseModel `bun:"table:items"`
	GroupItemID   uuid.UUID `bun:"group_item_id,type:uuid,notnull"`
	Name          string    `bun:"name,notnull"`
	Description   string    `bun:"description"`
	Observation   string    `bun:"observation"`
	Price         float64   `bun:"price,notnull"`
	Quantity      float64   `bun:"quantity,notnull"`
}
