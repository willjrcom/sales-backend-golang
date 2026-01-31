package companyusecases

import (
	"context"
	"time"

	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

func (s *Service) GetSubscriptionStatus(ctx context.Context) (*companydto.SubscriptionStatusDTO, error) {
	// Get company for current_plan and subscription_expires_at
	company, err := s.r.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	dto := &companydto.SubscriptionStatusDTO{
		CurrentPlan:      string(company.CurrentPlan),
		CanCancelRenewal: false,
		Periodicity:      "MONTHLY",
	}

	// Get active subscription - single source of truth
	activeSub, err := s.r.GetActiveSubscription(ctx, company.ID)
	if err == nil && activeSub != nil {
		// Expiration date and days remaining from active subscription
		expiresAt := activeSub.EndDate.Format(time.RFC3339)
		dto.ExpiresAt = &expiresAt

		daysRemaining := max(int(time.Until(activeSub.EndDate).Hours()/24), 0)
		dto.DaysRemaining = &daysRemaining

		// Can cancel if subscription hasn't been cancelled yet
		dto.CanCancelRenewal = !activeSub.IsCanceled

		// Get periodicity from linked payment
		if activeSub.PaymentID != nil {
			payment, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, *activeSub.PaymentID)
			if err == nil && payment != nil {
				switch payment.Months {
				case 6:
					dto.Periodicity = "SEMIANNUAL"
				case 12:
					dto.Periodicity = "ANNUAL"
				}
			}
		}
	}

	// Check for upcoming (future) subscription
	upcoming, err := s.r.GetUpcomingSubscription(ctx, company.ID)
	if err == nil && upcoming != nil {
		planType := string(upcoming.PlanType)
		startAt := upcoming.StartDate.Format(time.RFC3339)
		dto.UpcomingPlan = &planType
		dto.UpcomingStartAt = &startAt
	}

	return dto, nil
}
