package companyentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PlanType string

const (
	PlanTypeFree         PlanType = "free"
	PlanTypeBasic        PlanType = "basic"
	PlanTypeIntermediate PlanType = "intermediate"
	PlanTypeAdvanced     PlanType = "advanced"
)

type CompanySubscription struct {
	entity.Entity
	CompanyID uuid.UUID
	PaymentID *uuid.UUID // Can be nil for manual grants
	PlanType  PlanType
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

func NewCompanySubscription(companyID uuid.UUID, planType PlanType, startDate, endDate time.Time) *CompanySubscription {
	return &CompanySubscription{
		Entity:    entity.NewEntity(),
		CompanyID: companyID,
		PlanType:  planType,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}
}

func (c *CompanySubscription) IsExpired() bool {
	return time.Now().After(c.EndDate)
}
