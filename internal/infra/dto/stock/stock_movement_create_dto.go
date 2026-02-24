package stockdto

import (
	"github.com/shopspring/decimal"
)

// StockMovementCreateDTO representa o DTO para criar movimento
type StockMovementCreateDTO struct {
	Reason     string          `json:"reason"`
	Quantity   decimal.Decimal `json:"quantity,omitempty"`
	Price      decimal.Decimal `json:"price,omitempty"`
	TotalPrice decimal.Decimal `json:"total_price,omitempty"`
}
