package productentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type CategoryProduct struct {
	entity.Entity
	Name     string
	Sizes    []string
	Products []*Product `bun:"rel:has-many"`
}
