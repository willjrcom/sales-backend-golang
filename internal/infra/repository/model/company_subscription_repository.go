package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanySubscriptionRepository interface {
	// Subscriptions
	CreateSubscription(ctx context.Context, subscription *CompanySubscription) error
	UpdateSubscription(ctx context.Context, subscription *CompanySubscription) error
	MarkSubscriptionAsCancelled(ctx context.Context, companyID uuid.UUID, externalReference string) error
	MarkSubscriptionAsActive(ctx context.Context, companyID uuid.UUID, externalReference string) error
	UpdateSubscriptionStatus(ctx context.Context, companyID uuid.UUID, status string, externalReference string) error
	GetActiveSubscription(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, error)
	GetLastPlan(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, error)
	GetByPreapprovalID(ctx context.Context, preapprovalID string) (*CompanySubscription, error)
	GetByExternalReference(ctx context.Context, externalReference string) (*CompanySubscription, error)
	UpdateCompanyPlans(ctx context.Context) error
}
