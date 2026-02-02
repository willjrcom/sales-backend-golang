package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	NewCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context, withoutRelations ...bool) (*Company, error)
	ListPublicCompanies(ctx context.Context) ([]Company, error)
	ListCompaniesForBilling(ctx context.Context) ([]Company, error)
	ListCompaniesByPaymentDueDay(ctx context.Context, day int) ([]Company, error)
	ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error)

	// Users
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
	GetCompanyUsers(ctx context.Context, page, perPage int) ([]User, int, error)
	UpdateBlockStatus(ctx context.Context, companyID uuid.UUID, isBlocked bool) error
}
