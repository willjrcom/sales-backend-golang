package productentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProcessRule struct {
	entity.Entity
	ProcessRuleCommonAttributes
}

type ProcessRuleCommonAttributes struct {
	Name              string
	Order             int8
	Description       string
	ImagePath         *string
	IdealTime         time.Duration
	ExperimentalError time.Duration
	CategoryID        uuid.UUID
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
