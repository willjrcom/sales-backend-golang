package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

// StockDTO representa o DTO de estoque
type StockDTO struct {
	ID           uuid.UUID                     `json:"id"`
	ProductID    uuid.UUID                     `json:"product_id"`
	Product      productcategorydto.ProductDTO `json:"product"`
	CurrentStock decimal.Decimal               `json:"current_stock"`
	MinStock     decimal.Decimal               `json:"min_stock"`
	MaxStock     decimal.Decimal               `json:"max_stock"`
	Unit         string                        `json:"unit"`
	IsActive     bool                          `json:"is_active"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// StockUpdateDTO representa o DTO para atualizar estoque

// FromDomain converte domain para DTO
func (s *StockDTO) FromDomain(stock *stockentity.Stock) {
	if stock == nil {
		return
	}

	productDTO := productcategorydto.ProductDTO{}
	productDTO.FromDomain(&stock.Product)

	*s = StockDTO{
		ID:           stock.ID,
		ProductID:    stock.ProductID,
		Product:      productDTO,
		CurrentStock: stock.CurrentStock,
		MinStock:     stock.MinStock,
		MaxStock:     stock.MaxStock,
		Unit:         stock.Unit,
		IsActive:     stock.IsActive,
		CreatedAt:    stock.CreatedAt,
		UpdatedAt:    stock.UpdatedAt,
	}
}
