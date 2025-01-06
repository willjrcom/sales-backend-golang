package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Item struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_items"`
	ItemCommonAttributes
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
