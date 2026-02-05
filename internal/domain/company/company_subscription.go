package companyentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type CompanySubscription struct {
	entity.Entity
	CompanyID         uuid.UUID
	PlanType          PlanType
	StartDate         time.Time
	EndDate           time.Time
	IsActive          bool
	IsCancelled       bool    // If true, renewal (Preapproval) was cancelled in MercadoPago
	PreapprovalID     *string // Unique identifier for the subscription contract
	ExternalReference *string // Reference for the subscription flow
	Status            string
}

func NewCompanySubscription(companyID uuid.UUID, planType PlanType, startDate, endDate time.Time) *CompanySubscription {
	ent := entity.NewEntity()
	return &CompanySubscription{
		Entity:            ent,
		CompanyID:         companyID,
		PlanType:          planType,
		StartDate:         startDate,
		EndDate:           endDate,
		IsActive:          true,
		IsCancelled:       false,
		PreapprovalID:     nil,
		ExternalReference: nil,
	}
}

func (c *CompanySubscription) IsExpired() bool {
	return time.Now().UTC().After(c.EndDate)
}
