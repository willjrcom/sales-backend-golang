package billingdto

import (
	"strings"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CreateSubscriptionCheckoutDTO struct {
	Plan      string `json:"plan"`      // "BASIC", "INTERMEDIATE", "ADVANCED"
	Frequency string `json:"frequency"` // "MONTHLY", "SEMIANNUAL", "ANNUAL"
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

func (c *CreateSubscriptionCheckoutDTO) ToFrequency() companyentity.Frequency {
	switch strings.ToUpper(c.Frequency) {
	case "SEMIANNUAL":
		return companyentity.FrequencySemiannual
	case "ANNUAL":
		return companyentity.FrequencyAnnual
	default:
		return companyentity.FrequencyMonthly
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
