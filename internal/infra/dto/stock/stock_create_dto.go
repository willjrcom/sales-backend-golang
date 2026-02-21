package stockdto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockCreateDTO representa o DTO para criar estoque
type StockCreateDTO struct {
	ProductID          uuid.UUID       `json:"product_id"`
	ProductVariationID *uuid.UUID      `json:"product_variation_id"`
	CurrentStock       decimal.Decimal `json:"current_stock"`
	MinStock           decimal.Decimal `json:"min_stock"`
	MaxStock           decimal.Decimal `json:"max_stock"`
	Unit               string          `json:"unit"`
	IsActive           bool            `json:"is_active"`
}

// ToDomain converte DTO para domain
func (s *StockCreateDTO) ToDomain() *stockentity.Stock {
	return stockentity.NewStock(
		s.ProductID,
		s.ProductVariationID,
		s.CurrentStock,
		s.MinStock,
		s.MaxStock,
		s.Unit,
	)
}
