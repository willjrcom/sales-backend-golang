package model

import (
	"github.com/uptrace/bun"
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
