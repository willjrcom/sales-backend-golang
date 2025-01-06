package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Place struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:places"`
	PlaceCommonAttributes
}

type PlaceCommonAttributes struct {
	Name        string          `bun:"name,notnull"`
	ImagePath   *string         `bun:"image_path"`
	IsAvailable bool            `bun:"is_available"`
	Tables      []PlaceToTables `bun:"rel:has-many,join:id=place_id"`
}

type PatchPlace struct {
	Name        *string `json:"name"`
	ImagePath   *string `json:"image_path"`
	IsAvailable *bool   `json:"is_available"`
}

type PlaceToTables struct {
	PlaceID uuid.UUID `bun:"type:uuid,pk"`
	Place   *Place    `bun:"rel:belongs-to,join:place_id=id"`
	TableID uuid.UUID `bun:"type:uuid,pk"`
	Table   *Table    `bun:"rel:belongs-to,join:table_id=id"`
	Column  int       `bun:"column:column,notnull"`
	Row     int       `bun:"column:row,notnull"`
}
