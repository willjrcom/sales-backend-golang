package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	fiscalinvoice "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_invoice"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type FiscalInvoice struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:fiscal_invoices"`

	CompanyID          uuid.UUID  `bun:"company_id,type:uuid,notnull"`
	OrderID            uuid.UUID  `bun:"order_id,type:uuid,notnull"`
	ChaveAcesso        string     `bun:"chave_acesso,unique"`
	Numero             int        `bun:"numero"`
	Serie              int        `bun:"serie"`
	Status             string     `bun:"status,notnull"`
	XMLPath            string     `bun:"xml_path"`
	PDFPath            string     `bun:"pdf_path"`
	Protocolo          string     `bun:"protocolo"`
	ErrorMessage       string     `bun:"error_message"`
	EmittedAt          *time.Time `bun:"emitted_at"`
	CancelledAt        *time.Time `bun:"cancelled_at"`
	CancellationReason string     `bun:"cancellation_reason"`
}

func (f *FiscalInvoice) FromDomain(invoice *fiscalinvoice.FiscalInvoice) {
	if invoice == nil {
		return
	}
	*f = FiscalInvoice{
		Entity:             entitymodel.FromDomain(invoice.Entity),
		CompanyID:          invoice.CompanyID,
		OrderID:            invoice.OrderID,
		ChaveAcesso:        invoice.ChaveAcesso,
		Numero:             invoice.Numero,
		Serie:              invoice.Serie,
		Status:             string(invoice.Status),
		XMLPath:            invoice.XMLPath,
		PDFPath:            invoice.PDFPath,
		Protocolo:          invoice.Protocolo,
		ErrorMessage:       invoice.ErrorMessage,
		CancellationReason: invoice.CancellationReason,
	}

	if invoice.IsAuthorized() && f.EmittedAt == nil {
		now := time.Now()
		f.EmittedAt = &now
	}

	if invoice.IsCancelled() && f.CancelledAt == nil {
		now := time.Now()
		f.CancelledAt = &now
	}
}

func (f *FiscalInvoice) ToDomain() *fiscalinvoice.FiscalInvoice {
	if f == nil {
		return nil
	}
	return &fiscalinvoice.FiscalInvoice{
		Entity:             f.Entity.ToDomain(),
		CompanyID:          f.CompanyID,
		OrderID:            f.OrderID,
		ChaveAcesso:        f.ChaveAcesso,
		Numero:             f.Numero,
		Serie:              f.Serie,
		Status:             fiscalinvoice.InvoiceStatus(f.Status),
		XMLPath:            f.XMLPath,
		PDFPath:            f.PDFPath,
		Protocolo:          f.Protocolo,
		ErrorMessage:       f.ErrorMessage,
		CancellationReason: f.CancellationReason,
	}
}

// FiscalInvoiceRepository defines the interface for fiscal invoice operations
type FiscalInvoiceRepository interface {
	Create(ctx context.Context, invoice *FiscalInvoice) error
	Update(ctx context.Context, invoice *FiscalInvoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*FiscalInvoice, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*FiscalInvoice, error)
	GetByChaveAcesso(ctx context.Context, chaveAcesso string) (*FiscalInvoice, error)
	List(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]*FiscalInvoice, int, error)
	GetNextNumero(ctx context.Context, companyID uuid.UUID, serie int) (int, error)
}
