package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StockReportDTO representa o DTO para relatórios de estoque
type StockReportDTO struct {
	ProductID    uuid.UUID       `json:"product_id"`
	ProductName  string          `json:"product_name"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"`
	MaxStock     decimal.Decimal `json:"max_stock"`
	Unit         string          `json:"unit"`
	StockLevel   decimal.Decimal `json:"stock_level"` // Porcentagem
	IsLowStock   bool            `json:"is_low_stock"`
	IsOutOfStock bool            `json:"is_out_of_stock"`
	TotalIn      decimal.Decimal `json:"total_in"`
	TotalOut     decimal.Decimal `json:"total_out"`
	TotalPrice   decimal.Decimal `json:"total_cost"`
}

// StockReportSummaryDTO representa o resumo do relatório de estoque
type StockReportSummaryDTO struct {
	TotalProducts     int             `json:"total_products"`
	TotalLowStock     int             `json:"total_low_stock"`
	TotalOutOfStock   int             `json:"total_out_of_stock"`
	TotalActiveAlerts int             `json:"total_active_alerts"`
	TotalStockValue   decimal.Decimal `json:"total_stock_value"`
}

// StockReportCompleteDTO representa o relatório completo de estoque
type StockReportCompleteDTO struct {
	Summary            StockReportSummaryDTO `json:"summary"`
	AllStocks          []StockDTO            `json:"all_stocks"`
	LowStockProducts   []StockDTO            `json:"low_stock_products"`
	OutOfStockProducts []StockDTO            `json:"out_of_stock_products"`
	ActiveAlerts       []StockAlertDTO       `json:"active_alerts"`
	GeneratedAt        time.Time             `json:"generated_at"`
}
