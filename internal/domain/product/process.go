package productentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Process struct {
	entity.Entity
	bun.BaseModel `bun:"table:processes"`
	ProcessCommonAttributes
}

type ProcessCommonAttributes struct {
	Name              string         `bun:"name,notnull" json:"name"`
	Order             int8           `bun:"order,notnull" json:"order"`
	IdealTime         *time.Duration `bun:"ideal_time,notnull" json:"ideal_time"`
	ExperimentalError *time.Duration `bun:"experimental_error,notnull" json:"experimental_error"`
	CategoryID        uuid.UUID      `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
}

type PatchProcess struct {
	Name              *string        `json:"name"`
	Order             *int8          `json:"order"`
	IdealTime         *time.Duration `json:"ideal_time"`
	ExperimentalError *time.Duration `json:"experimental_error"`
}
