package tableentity

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Table struct {
	entity.Entity
	bun.BaseModel `bun:"table:tables"`
	TableCommonAttributes
}

type TableCommonAttributes struct {
	Name        string                   `bun:"name,notnull" json:"name"`
	IsAvailable bool                     `bun:"is_available" json:"is_available"`
	Orders      []orderentity.TableOrder `bun:"rel:has-many,join:id=table_id" json:"orders"`
}

func (t *Table) LockTable() {
	t.IsAvailable = false
}

func (t *Table) UnlockTable() {
	t.IsAvailable = true
}
