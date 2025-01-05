package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductCategoryToAdditional struct {
	bun.BaseModel
	CategoryID           uuid.UUID
	Category             *ProductCategory
	AdditionalCategoryID uuid.UUID
	AdditionalCategory   *ProductCategory
}

type CategoryRelation struct {
	ID uuid.UUID
}
