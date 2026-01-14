package companyentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// CostType represents different types of usage costs
type CostType string

const (
	CostTypeSubscription       CostType = "subscription"
	CostTypeFiscalSubscription CostType = "fiscal_subscription"
	CostTypeNFCe               CostType = "nfce"
	CostTypeNFe                CostType = "nfe"
	CostTypeNFCeRefund         CostType = "nfce_refund"
	CostTypeNFeRefund          CostType = "nfe_refund"
)

var (
	// NFCeCost is the cost per NFC-e emission (R$ 0.10)
	NFCeCost = decimal.NewFromFloat(0.10)
	// MonthlyFiscalFee is the monthly subscription fee for fiscal invoice functionality (R$ 20.00)
	MonthlyFiscalFee = decimal.NewFromFloat(20.00)
)

// CompanyUsageCost represents a cost incurred by a company for using platform features
type CompanyUsageCost struct {
	entity.Entity
	CompanyID    uuid.UUID
	CostType     CostType
	Description  string
	Amount       decimal.Decimal
	ReferenceID  *uuid.UUID // Optional reference to the entity that generated this cost (e.g., fiscal_invoice.id)
	BillingMonth int        // 1-12
	BillingYear  int        // e.g., 2026
}

// NewCompanyUsageCost creates a new usage cost record
func NewCompanyUsageCost(companyID uuid.UUID, costType CostType, amount decimal.Decimal, description string, referenceID *uuid.UUID) *CompanyUsageCost {
	now := time.Now()
	return &CompanyUsageCost{
		Entity:       entity.NewEntity(),
		CompanyID:    companyID,
		CostType:     costType,
		Amount:       amount,
		Description:  description,
		ReferenceID:  referenceID,
		BillingMonth: int(now.Month()),
		BillingYear:  now.Year(),
	}
}

// NewNFCeCost creates a cost record for NFC-e emission
func NewNFCeCost(companyID uuid.UUID, invoiceID uuid.UUID, description string) *CompanyUsageCost {
	return NewCompanyUsageCost(companyID, CostTypeNFCe, NFCeCost, description, &invoiceID)
}

// NewSubscriptionCost creates a cost record for subscription payment
func NewSubscriptionCost(companyID uuid.UUID, amount decimal.Decimal, paymentID uuid.UUID, description string) *CompanyUsageCost {
	return NewCompanyUsageCost(companyID, CostTypeSubscription, amount, description, &paymentID)
}

// NewFiscalSubscriptionCost creates a cost record for monthly fiscal subscription fee
func NewFiscalSubscriptionCost(companyID uuid.UUID) *CompanyUsageCost {
	return NewCompanyUsageCost(companyID, CostTypeFiscalSubscription, MonthlyFiscalFee, "Taxa mensal de emissão de notas fiscais", nil)
}
