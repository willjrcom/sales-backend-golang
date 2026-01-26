package model

import (
	"context"

	"github.com/google/uuid"
)

// CompanyUsageCostRepository defines the interface for company usage cost operations
type CompanyUsageCostRepository interface {
	Create(ctx context.Context, cost *CompanyUsageCost) error
	GetByID(ctx context.Context, id uuid.UUID) (*CompanyUsageCost, error)
	GetMonthlyCosts(ctx context.Context, companyID uuid.UUID, month, year int) ([]*CompanyUsageCost, error)
	Update(ctx context.Context, cost *CompanyUsageCost) error
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*CompanyUsageCost, error)
	GetPendingCosts(ctx context.Context, companyID uuid.UUID) ([]*CompanyUsageCost, error)
	UpdateCostsPaymentID(ctx context.Context, costIDs []uuid.UUID, paymentID uuid.UUID) error
}
