package companyentity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// CostType represents different types of usage costs
type CostType string

const (
	CostTypeNFCe       CostType = "nfce"
	CostTypeNFe        CostType = "nfe"
	CostTypeNFCeRefund CostType = "nfce_refund"
	CostTypeNFeRefund  CostType = "nfe_refund"
)

var (
	// NFCeCost is the cost per NFC-e emission (R$ 0.10)
	NFCeCost = decimal.NewFromFloat(0.10)
	// MonthlyFiscalFee is the monthly subscription fee for fiscal invoice functionality (R$ 20.00)
	MonthlyFiscalFee = decimal.NewFromFloat(20.00)
)

// CostStatus represents the status of a cost
type CostStatus string

const (
	CostStatusPending          CostStatus = "PENDING"
	CostStatusPaid             CostStatus = "PAID"
	CostStatusOverdue          CostStatus = "OVERDUE"
	CostStatusWaived           CostStatus = "WAIVED"
	CostStatusPaymentGenerated CostStatus = "PAYMENT_GENERATED"
)

// CompanyUsageCost represents a cost incurred by a company for using platform features
type CompanyUsageCost struct {
	entity.Entity
	CompanyID   uuid.UUID
	CostType    CostType
	Description string
	Amount      decimal.Decimal
	Status      CostStatus
	ReferenceID *uuid.UUID // Optional reference to the entity that generated this cost (e.g., fiscal_invoice.id)
	PaymentID   *uuid.UUID // Optional reference to the payment that settled this cost
}

// NewPlanUsageCost creates a new usage cost record
func NewPlanUsageCost(companyID uuid.UUID, costType CostType, amount decimal.Decimal, description string, referenceID *uuid.UUID) *CompanyUsageCost {
	return &CompanyUsageCost{
		Entity:      entity.NewEntity(),
		CompanyID:   companyID,
		CostType:    costType,
		Amount:      amount,
		Status:      CostStatusPending,
		Description: description,
		ReferenceID: referenceID,
	}
}

func NewUsageCost(companyID uuid.UUID, costType CostType, amount decimal.Decimal, description string, referenceID *uuid.UUID) *CompanyUsageCost {
	return &CompanyUsageCost{
		CompanyID:   companyID,
		CostType:    costType,
		Amount:      amount,
		Status:      CostStatusPending,
		Description: description,
		ReferenceID: referenceID,
	}
}
