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
	CompanyID  uuid.UUID
	PaymentID  *uuid.UUID // Link to the payment that created this subscription
	PlanType   PlanType
	StartDate  time.Time
	EndDate    time.Time
	IsActive   bool
	IsCanceled bool // If true, renewal (Preapproval) was cancelled in MercadoPago
}

func NewCompanySubscription(companyID uuid.UUID, planType PlanType, startDate, endDate time.Time) *CompanySubscription {
	ent := entity.NewEntity()
	return &CompanySubscription{
		Entity:     ent,
		CompanyID:  companyID,
		PlanType:   planType,
		StartDate:  startDate,
		EndDate:    endDate,
		IsActive:   true,
		IsCanceled: false,
	}
}

func (c *CompanySubscription) IsExpired() bool {
	return time.Now().UTC().After(c.EndDate)
}
