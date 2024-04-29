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
	Name                       string        `bun:"name,notnull" json:"name"`
	Order                      int8          `bun:"order,notnull" json:"order"`
	IdealTime                  time.Duration `bun:"ideal_time,notnull" json:"ideal_time"`
	ExperimentalError          time.Duration `bun:"experimental_error,notnull" json:"experimental_error"`
	IdealTimeFormatted         string        `bun:"ideal_time_formatted,notnull" json:"ideal_time_formatted"`
	ExperimentalErrorFormatted string        `bun:"experimental_error_formatted,notnull" json:"experimental_error_formatted"`
	CategoryID                 uuid.UUID     `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
}

type PatchProcessRule struct {
	Name              *string       `json:"name"`
	Order             *int8         `json:"order"`
	IdealTime         time.Duration `json:"ideal_time"`
	ExperimentalError time.Duration `json:"experimental_error"`
}

func NewProcessRule(processCommonAttributes ProcessRuleCommonAttributes) *ProcessRule {
	processCommonAttributes.IdealTimeFormatted = processCommonAttributes.IdealTime.String()
	processCommonAttributes.ExperimentalErrorFormatted = processCommonAttributes.ExperimentalError.String()

	return &ProcessRule{
		Entity:                      entity.NewEntity(),
		ProcessRuleCommonAttributes: processCommonAttributes,
	}
}
