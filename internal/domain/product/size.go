package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Size struct {
	entity.Entity
	bun.BaseModel `bun:"table:sizes"`
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string    `bun:"name" json:"name"`
	Active     *bool     `bun:"active" json:"active"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Products   []Product `bun:"rel:has-many,join:id=size_id" json:"products"`
}

type PatchSize struct {
	Name   *string `json:"name"`
	Active *bool   `json:"active"`
}
