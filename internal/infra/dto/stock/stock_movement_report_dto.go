package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StockMovementReportDTO representa o DTO para relatório de movimentos
type StockMovementReportDTO struct {
	Date         time.Time       `json:"date"`
	ProductID    uuid.UUID       `json:"product_id"`
	ProductName  string          `json:"product_name"`
	Type         string          `json:"type"`
	Quantity     decimal.Decimal `json:"quantity"`
	Reason       string          `json:"reason"`
	OrderNumber  *int            `json:"order_number,omitempty"`
	EmployeeName *string         `json:"employee_name,omitempty"`
	Price        decimal.Decimal `json:"unit_cost"`
	TotalPrice   decimal.Decimal `json:"total_cost"`
}
