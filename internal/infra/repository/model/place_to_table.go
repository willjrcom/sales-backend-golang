package model

import (
	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type PlaceToTables struct {
	PlaceID uuid.UUID `bun:"type:uuid,pk"`
	Place   *Place    `bun:"rel:belongs-to,join:place_id=id"`
	TableID uuid.UUID `bun:"type:uuid,pk"`
	Table   *Table    `bun:"rel:belongs-to,join:table_id=id"`
	Column  int       `bun:"column:column,notnull"`
	Row     int       `bun:"column:row,notnull"`
}

func (p *PlaceToTables) FromDomain(placeToTables *tableentity.PlaceToTables) {
	*p = PlaceToTables{
		PlaceID: placeToTables.PlaceID,
		Place:   &Place{},
		TableID: placeToTables.TableID,
		Table:   &Table{},
		Column:  placeToTables.Column,
		Row:     placeToTables.Row,
	}
}

func (p *PlaceToTables) ToDomain() *tableentity.PlaceToTables {
	if p == nil {
		return nil
	}
	return &tableentity.PlaceToTables{
		PlaceID: p.PlaceID,
		Place:   p.Place.ToDomain(),
		TableID: p.TableID,
		Table:   p.Table.ToDomain(),
		Column:  p.Column,
		Row:     p.Row,
	}
}
