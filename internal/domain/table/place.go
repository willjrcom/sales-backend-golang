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

type TableRelation struct {
	TableID uuid.UUID `bun:"column:table_id,type:uuid,notnull" json:"table_id"`
	Table   *Table    `bun:"rel:belongs-to" json:"table"`
}

type PlaceCommonAttributes struct {
	Tables [][]TableRelation `bun:"rel:has-many,join:id=place_id" json:"tables"`
}
