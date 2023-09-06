package productentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Product struct {
	entity.Entity
	Code        string
	Name        string
	Description string
	Size        string
	Price       float64
	Cost        float64
	CategoryID  uuid.UUID `bun:"rel:belongs-to`
	Category    *CategoryProduct
	IsAvailable bool
}
