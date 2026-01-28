package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CompanyPaymentRepository interface {
	CreateCompanyPayment(ctx context.Context, payment *CompanyPayment) error
	UpdateCompanyPayment(ctx context.Context, payment *CompanyPayment) error
	GetCompanyPaymentByID(ctx context.Context, id uuid.UUID) (*CompanyPayment, error)
	ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]CompanyPayment, int, error)
	ListOverduePayments(ctx context.Context, cutoffDate time.Time) ([]CompanyPayment, error)
	ListPendingMandatoryPayments(ctx context.Context, companyID uuid.UUID) ([]CompanyPayment, error)
	ListOverduePaymentsByCompany(ctx context.Context, companyID uuid.UUID, cutoffDate time.Time) ([]CompanyPayment, error)
	ListExpiredOptionalPayments(ctx context.Context) ([]CompanyPayment, error)
}
