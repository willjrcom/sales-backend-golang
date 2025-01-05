package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Table struct {
	entity.Entity
	bun.BaseModel `bun:"table:tables"`
	TableCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type TableCommonAttributes struct {
	Name        string                   `bun:"name,notnull"`
	IsAvailable bool                     `bun:"is_available"`
	Orders      []orderentity.OrderTable `bun:"rel:has-many,join:id=table_id"`
}
