package tableentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Place struct {
	entity.Entity
	bun.BaseModel `bun:"table:places"`
	PlaceCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type PlaceCommonAttributes struct {
	Name        string          `bun:"name,notnull" json:"name"`
	ImagePath   *string         `bun:"image_path" json:"image_path"`
	IsAvailable bool            `bun:"is_available" json:"is_available"`
	Tables      []PlaceToTables `bun:"rel:has-many,join:id=place_id" json:"tables,omitempty"`
}

type PatchPlace struct {
	Name        *string `json:"name"`
	ImagePath   *string `json:"image_path"`
	IsAvailable *bool   `json:"is_available"`
}

type PlaceToTables struct {
	PlaceID uuid.UUID `bun:"type:uuid,pk" json:"place_id,omitempty"`
	Place   *Place    `bun:"rel:belongs-to,join:place_id=id" json:"place,omitempty"`
	TableID uuid.UUID `bun:"type:uuid,pk" json:"table_id,omitempty"`
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
