package productentity

import (
	"github.com/google/uuid"
)

type ProductCategoryToComplement struct {
	CategoryID           uuid.UUID
	Category             *ProductCategory
	ComplementCategoryID uuid.UUID
	ComplementCategory   *ProductCategory
}
