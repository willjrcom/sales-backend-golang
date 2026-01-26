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
	ProviderPaymentID string          `json:"provider_payment_id"`
	Status            string          `json:"status"`
	Currency          string          `json:"currency"`
	Amount            decimal.Decimal `json:"amount"`
	Months            int             `json:"months"`
	PaidAt            *time.Time      `json:"paid_at,omitempty"`
	ExternalReference string          `json:"external_reference,omitempty"`
	PaymentURL        string          `json:"payment_url,omitempty"`
}

// FromDomain maps a domain subscription payment to a DTO representation.
func (c *CompanyPaymentDTO) FromDomain(payment *companyentity.SubscriptionPayment) {
	if payment == nil {
		return
	}

	*c = CompanyPaymentDTO{
		ID:                payment.ID,
		Provider:          payment.Provider,
		ProviderPaymentID: payment.ProviderPaymentID,
		Status:            payment.Status,
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		Months:            payment.Months,
		PaidAt:            payment.PaidAt,
		ExternalReference: payment.ExternalReference,
		PaymentURL:        payment.PaymentURL,
	}
}
