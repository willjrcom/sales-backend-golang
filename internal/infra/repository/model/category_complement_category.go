package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductCategoryToComplement struct {
	bun.BaseModel        `bun:"table:product_category_to_complement"`
	CategoryID           uuid.UUID        `bun:"column:category_id,type:uuid,pk"`
	Category             *ProductCategory `bun:"rel:belongs-to,join:category_id=id"`
	ComplementCategoryID uuid.UUID        `bun:"column:complement_category_id,type:uuid,pk"`
	ComplementCategory   *ProductCategory `bun:"rel:belongs-to,join:complement_category_id=id"`
}
