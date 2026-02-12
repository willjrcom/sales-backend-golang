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
	Name        string
	Order       int8
	Description string
	ImagePath   *string
	IdealTime   time.Duration
	CategoryID  uuid.UUID
	Category    *ProductCategory
	IsActive    bool
}

func NewProcessRule(processCommonAttributes ProcessRuleCommonAttributes) *ProcessRule {
	return &ProcessRule{
		Entity:                      entity.NewEntity(),
		ProcessRuleCommonAttributes: processCommonAttributes,
	}
}
