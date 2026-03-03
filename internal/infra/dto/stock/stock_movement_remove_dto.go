package stockdto

import (
	"github.com/shopspring/decimal"
)

// StockMovementRemoveDTO representa o DTO para remover movimento
type StockMovementRemoveDTO struct {
	Reason   string          `json:"reason"`
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price,omitempty"` // Opcional para saída manual
}
