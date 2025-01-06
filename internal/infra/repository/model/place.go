package model

import (
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
