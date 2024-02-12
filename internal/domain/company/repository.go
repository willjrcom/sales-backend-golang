package companyentity

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	NewCompany(ctx context.Context, company *Company) (uuid.UUID, error)
	GetCompany(ctx context.Context) (*Company, error)
	AddUser(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) error
}
