package productentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type CategoryProduct struct {
	entity.Entity
	Name     string     `bun:"name"`
	Sizes    []string   `bun:"sizes"`
	Products []*Product `bun:"products,rel:has-many,join:id=id"`
}
