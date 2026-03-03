package stockdto

import (
	"time"

	"github.com/shopspring/decimal"
)

// StockMovementCreateDTO representa o DTO para criar movimento
type StockMovementCreateDTO struct {
	Reason    string          `json:"reason"`
	Quantity  decimal.Decimal `json:"quantity"`
	Price     decimal.Decimal `json:"price"` // Preço de custo na entrada
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
}
