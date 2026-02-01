package model

import (
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
	AccessKey          string     `bun:"access_key,unique"` // ChaveAcesso
	Number             int        `bun:"number"`            // Numero
	Series             int        `bun:"series"`            // Serie
	Status             string     `bun:"status,notnull"`
	XMLPath            string     `bun:"xml_path"`
	PDFPath            string     `bun:"pdf_path"`
	Protocol           string     `bun:"protocol"` // Protocolo
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
		AccessKey:          invoice.AccessKey,
		Number:             invoice.Number,
		Series:             invoice.Series,
		Status:             string(invoice.Status),
		XMLPath:            invoice.XMLPath,
		PDFPath:            invoice.PDFPath,
		Protocol:           invoice.Protocol,
		ErrorMessage:       invoice.ErrorMessage,
		CancellationReason: invoice.CancellationReason,
	}

	if invoice.IsAuthorized() && f.EmittedAt == nil {
		now := time.Now().UTC()
		f.EmittedAt = &now
	}

	if invoice.IsCancelled() && f.CancelledAt == nil {
		now := time.Now().UTC()
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
		AccessKey:          f.AccessKey,
		Number:             f.Number,
		Series:             f.Series,
		Status:             fiscalinvoice.InvoiceStatus(f.Status),
		XMLPath:            f.XMLPath,
		PDFPath:            f.PDFPath,
		Protocol:           f.Protocol,
		ErrorMessage:       f.ErrorMessage,
		CancellationReason: f.CancellationReason,
	}
}
