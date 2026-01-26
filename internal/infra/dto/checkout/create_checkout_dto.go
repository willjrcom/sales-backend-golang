package billingdto

import (
	"github.com/google/uuid"
	domainbilling "github.com/willjrcom/sales-backend-go/internal/domain/checkout"
)

type CreateCheckoutDTO struct {
	CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Plan        string    `json:"plan"`        // "BASIC", "INTERMEDIATE", "ENTERPRISE"
	Periodicity string    `json:"periodicity"` // "MONTHLY", "SEMIANNUAL", "ANNUAL"
}

type CheckoutResponseDTO struct {
	PaymentID   string `json:"payment_id"`
	CheckoutUrl string `json:"checkout_url"`
}

func (c *CreateCheckoutDTO) ToPlanType() domainbilling.PlanType {
	switch c.Plan {
	case "INTERMEDIATE":
		return domainbilling.PlanIntermediate
	case "ENTERPRISE":
		return domainbilling.PlanEnterprise
	default:
		return domainbilling.PlanBasic
	}
}

func (c *CreateCheckoutDTO) ToPeriodicity() domainbilling.Periodicity {
	switch c.Periodicity {
	case "SEMIANNUAL":
		return domainbilling.PeriodicitySemiannual
	case "ANNUAL":
		return domainbilling.PeriodicityAnnual
	default:
		return domainbilling.PeriodicityMonthly
	}
}
