package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	NewCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context) (*Company, error)
	GetCompanyOnlyByID(ctx context.Context, id uuid.UUID) (*Company, error)
	ListPublicCompanies(ctx context.Context) ([]Company, error)
	ListCompaniesForBilling(ctx context.Context) ([]Company, error)
	ListCompaniesByPaymentDueDay(ctx context.Context, day int) ([]Company, error)
	UpdateCompanySubscription(ctx context.Context, companyID uuid.UUID, schema string, expiresAt *time.Time, planType string) error
	ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error)

	// Users
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
	GetCompanyUsers(ctx context.Context, page, perPage int) ([]User, int, error)
	UpdateBlockStatus(ctx context.Context, companyID uuid.UUID, isBlocked bool) error

	// Subscriptions
	CreateSubscription(ctx context.Context, subscription *CompanySubscription) error
	UpdateSubscription(ctx context.Context, subscription *CompanySubscription) error
	MarkActiveSubscriptionAsCanceled(ctx context.Context, companyID uuid.UUID) error
	GetActiveSubscription(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, error)
	GetUpcomingSubscription(ctx context.Context, companyID uuid.UUID) (*CompanySubscription, error)
	UpdateCompanyPlans(ctx context.Context) error
}
