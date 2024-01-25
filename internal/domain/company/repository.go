package companyentity

import (
	"context"
)

type Repository interface {
	NewCompany(ctx context.Context, company *Company) error
	GetCompanyById(ctx context.Context, id string) (*Company, error)
	GetAllCompaniesBySchemaName(ctx context.Context, schemaName string) ([]Company, error)
}
