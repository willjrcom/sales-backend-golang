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
	PaidAt            time.Time       `bun:"paid_at,notnull"`
	ExternalReference string          `bun:"external_reference"`
	Description       string          `bun:"description"`
	RawPayload        []byte          `bun:"raw_payload,type:jsonb"`
}

func (c *CompanyPayment) FromDomain(payment *companyentity.SubscriptionPayment) {
	if payment == nil {
		return
	}
	*c = CompanyPayment{
		Entity:            entitymodel.FromDomain(payment.Entity),
		CompanyID:         payment.CompanyID,
		Provider:          payment.Provider,
		ProviderPaymentID: payment.ProviderPaymentID,
		Status:            payment.Status,
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		Months:            payment.Months,
		PaidAt:            payment.PaidAt,
		ExternalReference: payment.ExternalReference,
		RawPayload:        payment.RawPayload,
	}
}

func (c *CompanyPayment) ToDomain() *companyentity.SubscriptionPayment {
	if c == nil {
		return nil
	}
	return &companyentity.SubscriptionPayment{
		Entity:            c.Entity.ToDomain(),
		CompanyID:         c.CompanyID,
		Provider:          c.Provider,
		ProviderPaymentID: c.ProviderPaymentID,
		Status:            c.Status,
		Currency:          c.Currency,
		Amount:            c.Amount,
		Months:            c.Months,
		PaidAt:            c.PaidAt,
		ExternalReference: c.ExternalReference,
		RawPayload:        c.RawPayload,
	}
}
