package orderentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Table struct {
	entity.Entity
	TableCommonAttributes
}

type TableCommonAttributes struct {
	Name        string
	IsAvailable bool
	IsActive    bool
	Orders      []OrderTable
}

func (t *Table) LockTable() {
	t.IsAvailable = false
}

func (t *Table) UnlockTable() {
	t.IsAvailable = true
}
