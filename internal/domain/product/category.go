package productentity

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Category struct {
	entity.Entity
	bun.BaseModel `bun:"table:categories"`
	Name          string    `bun:"name"`
	Sizes         []Size    `bun:"rel:has-many,join:id=category_id"`
	Products      []Product `bun:"rel:has-many,join:id=category_id"`
}
