package stockdto

import (
	"github.com/shopspring/decimal"
)

// StockMovementAdjustDTO representa o DTO para criar movimento
type StockMovementAdjustDTO struct {
	NewStock decimal.Decimal `json:"new_stock"`
	Reason   string          `json:"reason"`
}
