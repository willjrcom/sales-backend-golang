package orderentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type TableOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:table_orders,alias:table_order"`
	TableOrderCommonAttributes
}

type TableOrderCommonAttributes struct {
	Name    string    `bun:"name,notnull" json:"name,omitempty"`
	Contact string    `bun:"contact,notnull" json:"contact,omitempty"`
	OrderID uuid.UUID `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
	TableID uuid.UUID `bun:"column:table_id,type:uuid,notnull" json:"table_id"`
}

func NewTable(tableOrderCommonAttributes TableOrderCommonAttributes) *TableOrder {
	return &TableOrder{
		Entity:                     entity.NewEntity(),
		TableOrderCommonAttributes: tableOrderCommonAttributes,
	}
}
