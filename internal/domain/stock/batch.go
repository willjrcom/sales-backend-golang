package stockentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// StockBatch representa um lote específico de estoque
type StockBatch struct {
	entity.Entity
	StockBatchCommonAttributes
}

type StockBatchCommonAttributes struct {
	StockID            uuid.UUID
	ProductVariationID uuid.UUID
	InitialQuantity    decimal.Decimal
	CurrentQuantity    decimal.Decimal
	CostPrice          decimal.Decimal
	ExpiresAt          *time.Time
}

func NewStockBatch(stockID uuid.UUID, variationID uuid.UUID, quantity, costPrice decimal.Decimal, expiresAt *time.Time) *StockBatch {
	return &StockBatch{
		Entity: entity.NewEntity(),
		StockBatchCommonAttributes: StockBatchCommonAttributes{
			StockID:            stockID,
			ProductVariationID: variationID,
			InitialQuantity:    quantity,
			CurrentQuantity:    quantity,
			CostPrice:          costPrice,
			ExpiresAt:          expiresAt,
		},
	}
}

func (sb *StockBatch) IsExpired() bool {
	if sb.ExpiresAt == nil {
		return false
	}
	return sb.ExpiresAt.Before(time.Now())
}

func (sb *StockBatch) HasStock() bool {
	return sb.CurrentQuantity.GreaterThan(decimal.Zero)
}

func (sb *StockBatch) CheckExpiration(daysThreshold int) AlertType {
	if sb.ExpiresAt == nil || !sb.HasStock() {
		return ""
	}

	threshold := time.Now().AddDate(0, 0, daysThreshold)
	if sb.ExpiresAt.Before(time.Now()) {
		return AlertTypeExpired
	}

	if sb.ExpiresAt.Before(threshold) {
		return AlertTypeNearExpiration
	}

	return ""
}
