package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// StockBatch model
type StockBatch struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_batches,alias:stock_batch"`
	StockBatchCommonAttributes
}

type StockBatchCommonAttributes struct {
	StockID            uuid.UUID        `bun:"stock_id,type:uuid,notnull"`
	ProductVariationID *uuid.UUID       `bun:"product_variation_id,type:uuid"`
	InitialQuantity    *decimal.Decimal `bun:"initial_quantity,type:decimal(10,3),notnull"`
	CurrentQuantity    *decimal.Decimal `bun:"current_quantity,type:decimal(10,3),notnull"`
	CostPrice          *decimal.Decimal `bun:"cost_price,type:decimal(10,2),notnull"`
	ExpiresAt          *time.Time       `bun:"expires_at"`
}

// FromDomain converte domain para model
func (sb *StockBatch) FromDomain(batch *stockentity.StockBatch) {
	if batch == nil {
		return
	}
	sb.ID = batch.ID
	sb.StockID = batch.StockID
	sb.ProductVariationID = nilIfZeroUUID(batch.ProductVariationID)
	sb.InitialQuantity = &batch.InitialQuantity
	sb.CurrentQuantity = &batch.CurrentQuantity
	sb.CostPrice = &batch.CostPrice
	sb.ExpiresAt = batch.ExpiresAt
	sb.CreatedAt = batch.CreatedAt
}

// ToDomain converte model para domain
func (sb *StockBatch) ToDomain() *stockentity.StockBatch {
	if sb == nil {
		return nil
	}
	return &stockentity.StockBatch{
		Entity: sb.Entity.ToDomain(),
		StockBatchCommonAttributes: stockentity.StockBatchCommonAttributes{
			StockID:            sb.StockID,
			ProductVariationID: derefUUID(sb.ProductVariationID),
			InitialQuantity:    sb.GetInitialQuantity(),
			CurrentQuantity:    sb.GetCurrentQuantity(),
			CostPrice:          sb.GetCostPrice(),
			ExpiresAt:          sb.ExpiresAt,
		},
	}
}

func (sb *StockBatch) GetInitialQuantity() decimal.Decimal {
	if sb.InitialQuantity == nil {
		return decimal.Zero
	}
	return *sb.InitialQuantity
}

func (sb *StockBatch) GetCurrentQuantity() decimal.Decimal {
	if sb.CurrentQuantity == nil {
		return decimal.Zero
	}
	return *sb.CurrentQuantity
}

func (sb *StockBatch) GetCostPrice() decimal.Decimal {
	if sb.CostPrice == nil {
		return decimal.Zero
	}
	return *sb.CostPrice
}
