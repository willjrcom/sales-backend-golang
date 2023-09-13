package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Size struct {
	entity.Entity
	bun.BaseModel `bun:"table:sizes"`
	Name          string    `bun:"name"`
	Active        bool      `bun:"active"`
	CategoryID    uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	Products      []Product `bun:"rel:has-many,join:id=size_id"`
}
