package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Quantity struct {
	entity.Entity
	bun.BaseModel `bun:"table:quantities"`
	Quantity      float64   `bun:"quantity,notnull"`
	CategoryID    uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
}
