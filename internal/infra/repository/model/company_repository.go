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
	GetCompanyByIDPublic(ctx context.Context, id uuid.UUID) (*Company, error)
	ListPublicCompanies(ctx context.Context) ([]Company, error)
	ListCompaniesForBilling(ctx context.Context) ([]Company, error)
	UpdateCompanySubscription(ctx context.Context, companyID uuid.UUID, schema string, expiresAt *time.Time, isBlocked bool) error
	ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error)
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
	GetCompanyUsers(ctx context.Context, page, perPage int) ([]User, int, error)
	CreateCompanyPayment(ctx context.Context, payment *CompanyPayment) error
	GetCompanyPaymentByProviderID(ctx context.Context, provider string, paymentID string) (*CompanyPayment, error)
}
