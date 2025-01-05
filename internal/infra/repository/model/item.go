package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Item struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_items"`
	ItemCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type ItemCommonAttributes struct {
	Name            string    `bun:"name,notnull"`
	Observation     string    `bun:"observation"`
	Price           float64   `bun:"price,notnull"`
	TotalPrice      float64   `bun:"total_price,notnull"`
	Size            string    `bun:"size,notnull"`
	Quantity        float64   `bun:"quantity,notnull"`
	GroupItemID     uuid.UUID `bun:"group_item_id,type:uuid"`
	CategoryID      uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	AdditionalItems []Item    `bun:"m2m:item_to_additional,join:Item=AdditionalItem"`
	RemovedItems    []string  `bun:"removed_items,type:jsonb"`
	ProductID       uuid.UUID `bun:"product_id,type:uuid"`
}
