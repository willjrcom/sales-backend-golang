package tableentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Place struct {
	entity.Entity
	PlaceCommonAttributes
}

type PlaceCommonAttributes struct {
	Name        string
	ImagePath   *string
	IsAvailable bool
	Tables      []PlaceToTables
}

type PlaceToTables struct {
	PlaceID uuid.UUID
	Place   *Place
	TableID uuid.UUID
	Table   *Table
	Column  int
	Row     int
}

func NewPlace(placeCommonAttributes PlaceCommonAttributes) *Place {
	return &Place{
		Entity:                entity.NewEntity(),
		PlaceCommonAttributes: placeCommonAttributes,
	}
}

func NewPlaceToTable(placeID, tableID uuid.UUID, column, row int) *PlaceToTables {
	return &PlaceToTables{
		PlaceID: placeID,
		TableID: tableID,
		Column:  column,
		Row:     row,
	}
}
