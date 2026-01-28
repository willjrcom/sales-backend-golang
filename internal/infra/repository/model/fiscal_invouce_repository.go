package model

import (
	"context"

	"github.com/google/uuid"
)

// FiscalInvoiceRepository defines the interface for fiscal invoice operations
type FiscalInvoiceRepository interface {
	Create(ctx context.Context, invoice *FiscalInvoice) error
	Update(ctx context.Context, invoice *FiscalInvoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*FiscalInvoice, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*FiscalInvoice, error)
	GetByAccessKey(ctx context.Context, accessKey string) (*FiscalInvoice, error)
	List(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]*FiscalInvoice, int, error)
	GetNextNumber(ctx context.Context, companyID uuid.UUID, series int) (int, error)
}
