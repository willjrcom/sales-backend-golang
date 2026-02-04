package companyentity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusApproved  PaymentStatus = "approved"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusRefused   PaymentStatus = "refused"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type CompanyPayment struct {
	entity.Entity
	CompanyID         uuid.UUID
	Provider          string
	ProviderPaymentID *string
	Status            PaymentStatus
	Currency          string
	Amount            decimal.Decimal
	Months            int
	PaidAt            *time.Time
	ExternalReference string
	PaymentURL        string
	ExpiresAt         *time.Time
	IsMandatory       bool
	Description       string
	PlanType          PlanType // "basic", "intermediate", "advanced", or empty for non-subscription payments
	RawPayload        json.RawMessage
}
