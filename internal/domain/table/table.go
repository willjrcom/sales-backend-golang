package tableentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Table struct {
	entity.Entity
	TableCommonAttributes
}

type TableCommonAttributes struct {
	Name        string
	IsAvailable bool
	Orders      []orderentity.OrderTable
}

type PatchTable struct {
	Name        *string
	IsAvailable *bool
}

func (t *Table) LockTable() {
	t.IsAvailable = false
}

func (t *Table) UnlockTable() {
	t.IsAvailable = true
}
