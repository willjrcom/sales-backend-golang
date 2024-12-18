package productentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProcessRule struct {
	entity.Entity
	bun.BaseModel `bun:"table:process_rules,alias:pr"`
	ProcessRuleCommonAttributes
}

type ProcessRuleCommonAttributes struct {
	Name              string        `bun:"name,notnull" json:"name"`
	Order             int8          `bun:"order,notnull" json:"order"`
	Description       string        `bun:"description" json:"description"`
	ImagePath         *string       `bun:"image_path" json:"image_path"`
	IdealTime         time.Duration `bun:"ideal_time,notnull" json:"ideal_time"`
	ExperimentalError time.Duration `bun:"experimental_error,notnull" json:"experimental_error"`
	CategoryID        uuid.UUID     `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	DeletedAt         time.Time     `bun:",soft_delete,nullzero"`
}

type PatchProcessRule struct {
	Name              *string `json:"name"`
	Order             *int8   `json:"order"`
	Description       *string `json:"description"`
	ImagePath         *string `json:"image_path"`
	IdealTime         *string `json:"ideal_time"`
	ExperimentalError *string `json:"experimental_error"`
}

func NewProcessRule(processCommonAttributes ProcessRuleCommonAttributes) *ProcessRule {
	return &ProcessRule{
		Entity:                      entity.NewEntity(),
		ProcessRuleCommonAttributes: processCommonAttributes,
	}
}
