package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Size struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:sizes"`
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string    `bun:"name"`
	IsActive   *bool     `bun:"is_active"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	Products   []Product `bun:"rel:has-many,join:id=size_id"`
}
