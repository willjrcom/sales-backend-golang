package model

import (
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Table struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:tables"`
	TableCommonAttributes
}

type TableCommonAttributes struct {
	Name        string       `bun:"name,notnull"`
	IsAvailable bool         `bun:"is_available"`
	Orders      []OrderTable `bun:"rel:has-many,join:id=table_id"`
}

func (t *Table) FromDomain(table *tableentity.Table) {
	if table == nil {
		return
	}
	*t = Table{
		Entity: entitymodel.FromDomain(table.Entity),
		TableCommonAttributes: TableCommonAttributes{
			Name:        table.Name,
			IsAvailable: table.IsAvailable,
			Orders:      []OrderTable{},
		},
	}

	for _, order := range table.Orders {
		ot := OrderTable{}
		ot.FromDomain(&order)
		t.Orders = append(t.Orders, ot)
	}
}

func (t *Table) ToDomain() *tableentity.Table {
	if t == nil {
		return nil
	}
	table := &tableentity.Table{
		Entity: t.Entity.ToDomain(),
		TableCommonAttributes: tableentity.TableCommonAttributes{
			Name:        t.Name,
			IsAvailable: t.IsAvailable,
			Orders:      []orderentity.OrderTable{},
		},
	}

	for _, order := range t.Orders {
		table.Orders = append(table.Orders, *order.ToDomain())
	}

	return table
}
