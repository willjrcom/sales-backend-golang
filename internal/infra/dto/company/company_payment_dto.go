package companydto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

// CompanyPaymentDTO represents a subscription payment made via Mercado Pago.
type CompanyPaymentDTO struct {
	ID                uuid.UUID       `json:"id"`
	Provider          string          `json:"provider"`
	ProviderPaymentID *string         `json:"provider_payment_id"`
	PreapprovalID     *string         `json:"preapproval_id"`
	Status            string          `json:"status"`
	Currency          string          `json:"currency"`
	Amount            decimal.Decimal `json:"amount"`
	Months            int             `json:"months"`
	PlanType          string          `json:"plan_type"`
	PaidAt            *time.Time      `json:"paid_at,omitempty"`
	ExternalReference string          `json:"external_reference,omitempty"`
	PaymentURL        string          `json:"payment_url,omitempty"`
	ExpiresAt         *time.Time      `json:"expires_at,omitempty"`
	IsMandatory       bool            `json:"is_mandatory"`
	CreatedAt         time.Time       `json:"created_at"`
	Description       string          `json:"description,omitempty"`
}

// FromDomain maps a domain subscription payment to a DTO representation.
func (c *CompanyPaymentDTO) FromDomain(payment *companyentity.CompanyPayment) {
	if payment == nil {
		return
	}

	*c = CompanyPaymentDTO{
		ID:                payment.ID,
		Provider:          payment.Provider,
		ProviderPaymentID: payment.ProviderPaymentID,
		PreapprovalID:     payment.PreapprovalID,
		Status:            string(payment.Status),
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		Months:            payment.Months,
		PlanType:          string(payment.PlanType),
		PaidAt:            payment.PaidAt,
		ExternalReference: payment.ExternalReference,
		PaymentURL:        payment.PaymentURL,
		ExpiresAt:         payment.ExpiresAt,
		IsMandatory:       payment.IsMandatory,
		CreatedAt:         payment.CreatedAt,
		Description:       payment.Description,
	}
}

// Extended DTO to include new fields if needed, or just add to the main one.
// Adding to main one is better.
