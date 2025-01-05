package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductCategoryToComplement struct {
	bun.BaseModel
	CategoryID           uuid.UUID
	Category             *ProductCategory
	ComplementCategoryID uuid.UUID
	ComplementCategory   *ProductCategory
}
