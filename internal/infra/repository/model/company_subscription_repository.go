package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanySubscriptionRepository interface {
	// Subscriptions
	CreateSubscription(ctx context.Context, subscription *CompanySubscription) error
	UpdateSubscription(ctx context.Context, subscription *CompanySubscription) error
	MarkActiveSubscriptionAsCanceled(ctx context.Context, companyID uuid.UUID) error
	GetActiveAndUpcomingSubscriptions(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, *CompanySubscription, error)
	UpdateCompanyPlans(ctx context.Context) error
}
