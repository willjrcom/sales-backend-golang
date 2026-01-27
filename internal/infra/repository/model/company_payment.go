package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyPayment struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:company_payments"`

	CompanyID         uuid.UUID       `bun:"company_id,type:uuid,notnull"`
	Provider          string          `bun:"provider,notnull"`
	ProviderPaymentID string          `bun:"provider_payment_id,notnull,unique"`
	Status            string          `bun:"status,notnull"`
	Currency          string          `bun:"currency,notnull"`
	Amount            decimal.Decimal `bun:"amount,type:decimal(12,2),notnull"`
	Months            int             `bun:"months,notnull"`
	PaidAt            *time.Time      `bun:"paid_at"`
	ExternalReference string          `bun:"external_reference"`
	PaymentURL        string          `bun:"payment_url"`
	ExpiresAt         *time.Time      `bun:"expires_at"`
	IsMandatory       bool            `bun:"is_mandatory,notnull,default:false"`
	Description       string          `bun:"description"`
	RawPayload        []byte          `bun:"raw_payload,type:jsonb"`
}

func (c *CompanyPayment) FromDomain(payment *companyentity.CompanyPayment) {
	if payment == nil {
		return
	}
	*c = CompanyPayment{
		Entity:            entitymodel.FromDomain(payment.Entity),
		CompanyID:         payment.CompanyID,
		Provider:          payment.Provider,
		ProviderPaymentID: payment.ProviderPaymentID,
		Status:            string(payment.Status),
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		Months:            payment.Months,
		PaidAt:            payment.PaidAt,
		ExternalReference: payment.ExternalReference,
		PaymentURL:        payment.PaymentURL,
		ExpiresAt:         payment.ExpiresAt,
		IsMandatory:       payment.IsMandatory,
		RawPayload:        payment.RawPayload,
	}
}

func (c *CompanyPayment) ToDomain() *companyentity.CompanyPayment {
	if c == nil {
		return nil
	}
	return &companyentity.CompanyPayment{
		Entity:            c.Entity.ToDomain(),
		CompanyID:         c.CompanyID,
		Provider:          c.Provider,
		ProviderPaymentID: c.ProviderPaymentID,
		Status:            companyentity.PaymentStatus(c.Status),
		Currency:          c.Currency,
		Amount:            c.Amount,
		Months:            c.Months,
		PaidAt:            c.PaidAt,
		ExternalReference: c.ExternalReference,
		PaymentURL:        c.PaymentURL,
		ExpiresAt:         c.ExpiresAt,
		IsMandatory:       c.IsMandatory,
		RawPayload:        c.RawPayload,
	}
}
