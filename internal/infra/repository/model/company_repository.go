package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	NewCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context) (*Company, error)
	ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error)
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
	GetCompanyUsers(ctx context.Context, offset, limit int) ([]User, int, error)
}
