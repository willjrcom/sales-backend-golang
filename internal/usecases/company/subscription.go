package companyusecases

import (
	"context"
	"time"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

func (s *Service) GetSubscriptionStatus(ctx context.Context) (*companydto.SubscriptionStatusDTO, error) {
	// Get company for current_plan and subscription_expires_at
	company, err := s.r.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	dto := &companydto.SubscriptionStatusDTO{
		CurrentPlan:      string(companyentity.PlanFree),
		CanCancelRenewal: false,
		Periodicity:      "MONTHLY",
	}

	// Get active and upcoming subscriptions - single call
	activeSub, err := s.companySubscriptionRepo.GetActiveSubscription(ctx, company.ID)
	if err == nil && activeSub != nil {
		// Expiration date and days remaining from active subscription
		expiresAt := activeSub.EndDate.Format(time.RFC3339)
		dto.ExpiresAt = &expiresAt

		daysRemaining := max(int(time.Until(activeSub.EndDate).Hours()/24), 0)
		dto.DaysRemaining = &daysRemaining

		// Can cancel if subscription hasn't been cancelled yet
		dto.CanCancelRenewal = !activeSub.IsCanceled

		monthsBetween := MonthsBetween(activeSub.StartDate, activeSub.EndDate)
		// Get periodicity from linked payment
		switch monthsBetween {
		case 6:
			dto.Periodicity = "SEMIANNUAL"
		case 12:
			dto.Periodicity = "ANNUAL"
		default:
			dto.Periodicity = "MONTHLY"
		}
	}

	return dto, nil
}

func MonthsBetween(a, b time.Time) int {
	// Garantir ordem (a <= b)
	if a.After(b) {
		a, b = b, a
	}

	years := b.Year() - a.Year()
	months := int(b.Month()) - int(a.Month())

	total := years*12 + months

	// Ajuste se ainda não completou o mês
	if b.Day() < a.Day() {
		total--
	}

	return total
}
