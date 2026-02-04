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
	GetCompanyPaymentByProviderID(ctx context.Context, providerPaymentID string) (*CompanyPayment, error)
	ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage, month, year int) ([]CompanyPayment, int, error)
	ListOverduePayments(ctx context.Context, cutoffDate time.Time) ([]CompanyPayment, error)
	ListOverduePaymentsByCompany(ctx context.Context, companyID uuid.UUID, cutoffDate time.Time) ([]CompanyPayment, error)
	ListPendingMandatoryPayments(ctx context.Context, companyID uuid.UUID) ([]CompanyPayment, error)
	GetCompanyPaymentByExternalReference(ctx context.Context, externalReference string) (*CompanyPayment, error)
	GetCompanyPaymentByExternalReferenceAndProviderID(ctx context.Context, externalReference string, providerPaymentID string) (*CompanyPayment, error)
	ListExpiredOptionalPayments(ctx context.Context) ([]CompanyPayment, error)
	ListByExternalReference(ctx context.Context, externalReference string) ([]CompanyPayment, error)

	// not used
	GetPendingPaymentByExternalReference(ctx context.Context, externalReference string) (*CompanyPayment, error)
}
