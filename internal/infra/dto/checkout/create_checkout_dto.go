package billingdto

import (
	"strings"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CreateSubscriptionCheckoutDTO struct {
	Plan        string `json:"plan"`        // "BASIC", "INTERMEDIATE", "ADVANCED"
	Periodicity string `json:"periodicity"` // "MONTHLY", "SEMIANNUAL", "ANNUAL"
}

type CheckoutResponseDTO struct {
	CheckoutUrl string `json:"checkout_url"`
}

func (c *CreateSubscriptionCheckoutDTO) ToPlanType() companyentity.PlanType {
	switch strings.ToUpper(c.Plan) {
	case "INTERMEDIATE":
		return companyentity.PlanIntermediate
	case "ADVANCED":
		return companyentity.PlanAdvanced
	default:
		return companyentity.PlanBasic
	}
}

func (c *CreateSubscriptionCheckoutDTO) ToPeriodicity() companyentity.Periodicity {
	switch strings.ToUpper(c.Periodicity) {
	case "SEMIANNUAL":
		return companyentity.PeriodicitySemiannual
	case "ANNUAL":
		return companyentity.PeriodicityAnnual
	default:
		return companyentity.PeriodicityMonthly
	}
}

type UpgradeSimulationDTO struct {
	TargetPlan     string  `json:"target_plan"`
	OldPlan        string  `json:"old_plan"`
	DaysRemaining  int     `json:"days_remaining"`
	UpgradeAmount  float64 `json:"upgrade_amount"`
	NewMonthlyCost float64 `json:"new_monthly_cost"`
	Frequency      int     `json:"frequency"`
}

type UpgradeCheckoutDTO struct {
	TargetPlan string `json:"target_plan"`
}
