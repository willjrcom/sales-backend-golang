package productentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Product struct {
	entity.Entity
	Code        string           `bun:"code"`
	Name        string           `bun:"name"`
	Description string           `bun:"description"`
	Size        string           `bun:"size"`
	Price       float64          `bun:"price"`
	Cost        float64          `bun:"cost"`
	CategoryID  uuid.UUID        `bun:"category_id,notnull"`
	Category    *CategoryProduct `bun:"rel:has-one"`
	IsAvailable bool             `bun:"is_available"`
}
