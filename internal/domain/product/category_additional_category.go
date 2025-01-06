package productentity

import (
	"github.com/google/uuid"
)

type ProductCategoryToAdditional struct {
	CategoryID           uuid.UUID
	Category             *ProductCategory
	AdditionalCategoryID uuid.UUID
	AdditionalCategory   *ProductCategory
}
