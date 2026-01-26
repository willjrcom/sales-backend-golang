package companydto

import "github.com/shopspring/decimal"

type MonthlyCostSummaryDTO struct {
	CompanyID   string                     `json:"company_id"`
	Month       int                        `json:"month"`
	Year        int                        `json:"year"`
	TotalAmount decimal.Decimal            `json:"total_amount"`
	CostsByType map[string]decimal.Decimal `json:"costs_by_type"`
	CostsCount  int                        `json:"costs_count"`
	OtherFee    decimal.Decimal            `json:"other_fee"`
	NFCeCosts   decimal.Decimal            `json:"nfce_costs"`
	NFCeCount   int                        `json:"nfce_count"`
}
