package stockdto

import (
	"github.com/shopspring/decimal"
)

// StockMovementRemoveDTO representa o DTO para criar movimento
type StockMovementRemoveDTO struct {
	Reason     string          `json:"reason"`
	Quantity   decimal.Decimal `json:"quantity,omitempty"`
	Price      decimal.Decimal `json:"price,omitempty"`
	TotalPrice decimal.Decimal `json:"total_price,omitempty"`
}
