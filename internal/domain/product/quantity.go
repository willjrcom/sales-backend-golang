package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Quantity struct {
	entity.Entity
	bun.BaseModel `bun:"table:quantities"`
	Name          string    `bun:"name"`
	CategoryID    uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
}
