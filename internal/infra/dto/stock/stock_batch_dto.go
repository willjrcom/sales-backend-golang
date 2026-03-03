package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockBatchDTO representa um lote de estoque para exibição no frontend
type StockBatchDTO struct {
	ID                 uuid.UUID       `json:"id"`
	StockID            uuid.UUID       `json:"stock_id"`
	ProductVariationID *uuid.UUID      `json:"product_variation_id,omitempty"`
	InitialQuantity    decimal.Decimal `json:"initial_quantity"`
	CurrentQuantity    decimal.Decimal `json:"current_quantity"`
	CostPrice          decimal.Decimal `json:"cost_price"`
	ExpiresAt          *time.Time      `json:"expires_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}

func (d *StockBatchDTO) FromDomain(b *stockentity.StockBatch) {
	d.ID = b.ID
	d.StockID = b.StockID
	if b.ProductVariationID != (uuid.UUID{}) {
		id := b.ProductVariationID
		d.ProductVariationID = &id
	}
	d.InitialQuantity = b.InitialQuantity
	d.CurrentQuantity = b.CurrentQuantity
	d.CostPrice = b.CostPrice
	d.ExpiresAt = b.ExpiresAt
	d.CreatedAt = b.CreatedAt
}
