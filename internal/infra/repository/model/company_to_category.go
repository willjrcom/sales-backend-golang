package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CompanyToCategory struct {
	bun.BaseModel `bun:"table:company_to_category"`
	CompanyID     uuid.UUID        `bun:"type:uuid,pk"`
	Company       *Company         `bun:"rel:belongs-to,join:company_id=id"`
	CategoryID    uuid.UUID        `bun:"type:uuid,pk"`
	Category      *CompanyCategory `bun:"rel:belongs-to,join:category_id=id"`
}
