package fiscalinvoice

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// InvoiceStatus represents the status of a fiscal invoice
type InvoiceStatus string

const (
	StatusPending    InvoiceStatus = "pending"
	StatusAuthorized InvoiceStatus = "authorized"
	StatusRejected   InvoiceStatus = "rejected"
	StatusCancelled  InvoiceStatus = "cancelled"
)

// FiscalInvoice represents a fiscal invoice (NFC-e or NF-e)
type FiscalInvoice struct {
	entity.Entity
	CompanyID          uuid.UUID
	OrderID            uuid.UUID
	AccessKey          string // ChaveAcesso (44-character access key)
	Number             int    // Numero
	Series             int    // Serie
	Status             InvoiceStatus
	XMLPath            string
	PDFPath            string
	Protocol           string // Protocolo
	ErrorMessage       string
	CancellationReason string
}

// NewFiscalInvoice creates a new fiscal invoice
func NewFiscalInvoice(companyID, orderID uuid.UUID, number, series int) *FiscalInvoice {
	return &FiscalInvoice{
		Entity:    entity.NewEntity(),
		CompanyID: companyID,
		OrderID:   orderID,
		Number:    number,
		Series:    series,
		Status:    StatusPending,
	}
}

// Authorize marks the invoice as authorized
func (f *FiscalInvoice) Authorize(accessKey, protocol, xmlPath, pdfPath string) {
	f.Status = StatusAuthorized
	f.AccessKey = accessKey
	f.Protocol = protocol
	f.XMLPath = xmlPath
	f.PDFPath = pdfPath
	f.ErrorMessage = ""
}

// Reject marks the invoice as rejected
func (f *FiscalInvoice) Reject(errorMessage string) {
	f.Status = StatusRejected
	f.ErrorMessage = errorMessage
}

// Cancel marks the invoice as cancelled
func (f *FiscalInvoice) Cancel(reason string) {
	f.Status = StatusCancelled
	f.CancellationReason = reason
}

// IsAuthorized checks if invoice is authorized
func (f *FiscalInvoice) IsAuthorized() bool {
	return f.Status == StatusAuthorized
}

// IsCancelled checks if invoice is cancelled
func (f *FiscalInvoice) IsCancelled() bool {
	return f.Status == StatusCancelled
}

// CanBeCancelled checks if invoice can be cancelled
func (f *FiscalInvoice) CanBeCancelled() bool {
	return f.Status == StatusAuthorized && f.AccessKey != ""
}
