package companydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyUsageCostDTO struct {
	ID           string          `json:"id"`
	CompanyID    string          `json:"company_id"`
	CostType     string          `json:"cost_type"`
	Description  string          `json:"description"`
	Amount       decimal.Decimal `json:"amount"`
	ReferenceID  *string         `json:"reference_id,omitempty"`
	BillingMonth int             `json:"billing_month"`
	BillingYear  int             `json:"billing_year"`
	CreatedAt    string          `json:"created_at"`
}

func (dto *CompanyUsageCostDTO) FromDomain(cost *companyentity.CompanyUsageCost) {
	if cost == nil {
		return
	}
	dto.ID = cost.ID.String()
	dto.CompanyID = cost.CompanyID.String()
	dto.CostType = string(cost.CostType)
	dto.Description = cost.Description
	dto.Amount = cost.Amount
	dto.BillingMonth = cost.BillingMonth
	dto.BillingYear = cost.BillingYear
	dto.CreatedAt = cost.CreatedAt.Format("2006-01-02T15:04:05Z07:00")

	if cost.ReferenceID != nil {
		refID := cost.ReferenceID.String()
		dto.ReferenceID = &refID
	}
}

type MonthlyCostSummaryDTO struct {
	CompanyID       string                     `json:"company_id"`
	Month           int                        `json:"month"`
	Year            int                        `json:"year"`
	TotalAmount     decimal.Decimal            `json:"total_amount"`
	CostsByType     map[string]decimal.Decimal `json:"costs_by_type"`
	CostsCount      int                        `json:"costs_count"`
	SubscriptionFee decimal.Decimal            `json:"subscription_fee"`
	NFCeCosts       decimal.Decimal            `json:"nfce_costs"`
	NFCeCount       int                        `json:"nfce_count"`
}

type CostBreakdownDTO struct {
	CompanyID   string                `json:"company_id"`
	Month       int                   `json:"month"`
	Year        int                   `json:"year"`
	Costs       []CompanyUsageCostDTO `json:"costs"`
	TotalAmount string                `json:"total_amount"`
}

type GetMonthlyCostsRequestDTO struct {
	Month int `json:"month" validate:"required,min=1,max=12"`
	Year  int `json:"year" validate:"required,min=2020"`
}

func (dto *GetMonthlyCostsRequestDTO) Validate() error {
	if dto.Month < 1 || dto.Month > 12 {
		return companyentity.ErrInvalidMonth
	}
	if dto.Year < 2020 {
		return companyentity.ErrInvalidYear
	}
	return nil
}

type RegisterUsageCostRequestDTO struct {
	CostType    string          `json:"cost_type" validate:"required"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	ReferenceID *uuid.UUID      `json:"reference_id,omitempty"`
}

// NextInvoicePreviewDTO represents the preview of the next billing invoice
type NextInvoicePreviewDTO struct {
	CompanyID         string              `json:"company_id"`
	NextBillingDate   string              `json:"next_billing_date"`
	EnabledServices   []EnabledServiceDTO `json:"enabled_services"`
	EstimatedTotal    decimal.Decimal     `json:"estimated_total"`
	CurrentMonthUsage decimal.Decimal     `json:"current_month_usage"`
	NFCeCount         int                 `json:"nfce_count"`
}

// EnabledServiceDTO represents a billable service and its status
type EnabledServiceDTO struct {
	Name        string          `json:"name"`
	Enabled     bool            `json:"enabled"`
	FixedCost   decimal.Decimal `json:"fixed_cost"`
	UsageCost   decimal.Decimal `json:"usage_cost"`
	UnitCost    decimal.Decimal `json:"unit_cost,omitempty"`
	UsageCount  int             `json:"usage_count,omitempty"`
	Description string          `json:"description"`
}
