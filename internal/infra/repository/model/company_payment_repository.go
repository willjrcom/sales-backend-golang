package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanyPaymentRepository interface {
	CreateCompanyPayment(ctx context.Context, payment *CompanyPayment) error
	UpdateCompanyPayment(ctx context.Context, payment *CompanyPayment) error
	GetCompanyPaymentByID(ctx context.Context, id uuid.UUID) (*CompanyPayment, error)
	ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]CompanyPayment, int, error)
}
