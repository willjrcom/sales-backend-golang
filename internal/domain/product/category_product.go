package productentity

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type CategoryProduct struct {
	entity.Entity
	bun.BaseModel `bun:"table:category_products"`
	Name          string    `bun:"name"`
	Sizes         []string  `bun:"sizes"`
	Products      []Product `bun:"rel:has-many,join:id=category_id"`
}
