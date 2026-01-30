package companyusecases

import (
	"context"
	"database/sql"
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
		CurrentPlan: string(company.CurrentPlan),
	}

	// If company has an expiration date, calculate days remaining
	if company.SubscriptionExpiresAt != nil {
		expiresAt := company.SubscriptionExpiresAt.Format(time.RFC3339)
		dto.ExpiresAt = &expiresAt

		daysRemaining := int(time.Until(*company.SubscriptionExpiresAt).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}
		dto.DaysRemaining = &daysRemaining
	}

	// Check for upcoming (future) subscriptions
	upcoming, err := s.r.GetUpcomingSubscription(ctx, company.ID)
	if err == nil && upcoming != nil {
		planType := string(upcoming.PlanType)
		startAt := upcoming.StartDate.Format(time.RFC3339)
		dto.UpcomingPlan = &planType
		dto.UpcomingStartAt = &startAt
	} else if err != nil && err != sql.ErrNoRows {
		// Log error but don't fail the entire request
		// upcoming subscription is optional information
	}

	// Check if company has an active subscription that can be cancelled
	// This means finding a pending payment with SUB:<companyID>: reference
	dto.CanCancelRenewal = false
	if company.ID.String() != "" {
		externalRef := "SUB:" + company.ID.String() + ":"
		payment, err := s.companyPaymentRepo.GetPendingPaymentByExternalReference(ctx, externalRef)
		if err == nil && payment != nil {
			dto.CanCancelRenewal = true
		}
	}

	return dto, nil
}
