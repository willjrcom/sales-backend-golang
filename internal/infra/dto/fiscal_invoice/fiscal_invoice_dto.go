package fiscalinvoicedto

import (
	"github.com/google/uuid"
	fiscalinvoice "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_invoice"
)

type FiscalInvoiceDTO struct {
	ID                 string `json:"id"`
	CompanyID          string `json:"company_id"`
	OrderID            string `json:"order_id"`
	ChaveAcesso        string `json:"chave_acesso,omitempty"`
	Numero             int    `json:"numero"`
	Serie              int    `json:"serie"`
	Status             string `json:"status"`
	XMLPath            string `json:"xml_path,omitempty"`
	PDFPath            string `json:"pdf_path,omitempty"`
	Protocolo          string `json:"protocolo,omitempty"`
	ErrorMessage       string `json:"error_message,omitempty"`
	CancellationReason string `json:"cancellation_reason,omitempty"`
	CreatedAt          string `json:"created_at"`
}

func (dto *FiscalInvoiceDTO) FromDomain(invoice *fiscalinvoice.FiscalInvoice) {
	if invoice == nil {
		return
	}
	dto.ID = invoice.ID.String()
	dto.CompanyID = invoice.CompanyID.String()
	dto.OrderID = invoice.OrderID.String()
	dto.ChaveAcesso = invoice.ChaveAcesso
	dto.Numero = invoice.Numero
	dto.Serie = invoice.Serie
	dto.Status = string(invoice.Status)
	dto.XMLPath = invoice.XMLPath
	dto.PDFPath = invoice.PDFPath
	dto.Protocolo = invoice.Protocolo
	dto.ErrorMessage = invoice.ErrorMessage
	dto.CancellationReason = invoice.CancellationReason
	dto.CreatedAt = invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
}

type EmitirNFCeRequestDTO struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

type CancelarNFCeRequestDTO struct {
	Justificativa string `json:"justificativa" validate:"required,min=15"`
}

type ListFiscalInvoicesRequestDTO struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
