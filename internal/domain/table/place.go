package tableentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Place struct {
	entity.Entity
	bun.BaseModel `bun:"table:places"`
	PlaceCommonAttributes
}

type PlaceCommonAttributes struct {
	Name        string  `bun:"name,notnull" json:"name"`
	ImagePath   *string `bun:"image_path" json:"image_path"`
	IsAvailable bool    `bun:"is_available" json:"is_available"`
	Tables      []Table `bun:"m2m:place_to_tables,join:Place=Table" json:"place_tables,omitempty"`
}

type PlaceToTables struct {
	PlaceID uuid.UUID `bun:"type:uuid,pk"`
	Place   *Place    `bun:"rel:belongs-to,join:place_id=id" json:"place,omitempty"`
	TableID uuid.UUID `bun:"type:uuid,pk"`
	Table   *Table    `bun:"rel:belongs-to,join:table_id=id" json:"table,omitempty"`
	Column  int       `bun:"column:column,notnull" json:"column"`
	Row     int       `bun:"column:row,notnull" json:"row"`
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
