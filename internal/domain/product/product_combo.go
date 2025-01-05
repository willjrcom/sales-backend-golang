package productentity

import (
	"github.com/google/uuid"
)

type ProductToCombo struct {
	ProductID uuid.UUID
	Product   *Product
}
