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
	MarkInactiveSubscriptionAsActive(ctx context.Context, companyID uuid.UUID) error
	GetActiveSubscription(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, error)
	GetByPreapprovalID(ctx context.Context, preapprovalID string) (*CompanySubscription, error)
	GetByExternalReference(ctx context.Context, externalReference string) (*CompanySubscription, error)
	UpdateCompanyPlans(ctx context.Context) error
}
