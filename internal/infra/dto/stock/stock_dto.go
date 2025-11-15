package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockDTO representa o DTO de estoque
type StockDTO struct {
	ID           uuid.UUID       `json:"id"`
	ProductID    uuid.UUID       `json:"product_id"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"`
	MaxStock     decimal.Decimal `json:"max_stock"`
	Unit         string          `json:"unit"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// StockUpdateDTO representa o DTO para atualizar estoque

// FromDomain converte domain para DTO
func (s *StockDTO) FromDomain(stock *stockentity.Stock) {
	if stock == nil {
		return
	}

	*s = StockDTO{
		ID:           stock.ID,
		ProductID:    stock.ProductID,
		CurrentStock: stock.CurrentStock,
		MinStock:     stock.MinStock,
		MaxStock:     stock.MaxStock,
		Unit:         stock.Unit,
		IsActive:     stock.IsActive,
		CreatedAt:    stock.CreatedAt,
		UpdatedAt:    stock.UpdatedAt,
	}
}
