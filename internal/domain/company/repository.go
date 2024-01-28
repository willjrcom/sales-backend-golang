package companyentity

import (
	"context"
)

type Repository interface {
	NewCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context) (*Company, error)
}
