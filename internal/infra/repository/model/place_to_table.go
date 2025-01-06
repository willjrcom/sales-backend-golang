package model

import "github.com/google/uuid"

type PlaceToTables struct {
	PlaceID uuid.UUID `bun:"type:uuid,pk"`
	Place   *Place    `bun:"rel:belongs-to,join:place_id=id"`
	TableID uuid.UUID `bun:"type:uuid,pk"`
	Table   *Table    `bun:"rel:belongs-to,join:table_id=id"`
	Column  int       `bun:"column:column,notnull"`
	Row     int       `bun:"column:row,notnull"`
}
