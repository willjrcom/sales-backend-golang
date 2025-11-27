package companyentity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type SubscriptionPayment struct {
	entity.Entity
	CompanyID         uuid.UUID
	Provider          string
	ProviderPaymentID string
	Status            string
	Currency          string
	Amount            decimal.Decimal
	Months            int
	PaidAt            time.Time
	ExternalReference string
	RawPayload        json.RawMessage
}
