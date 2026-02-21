package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// Stock model
type Stock struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stocks,alias:stock"`
	StockCommonAttributes
}

type StockCommonAttributes struct {
	ProductID          uuid.UUID       `bun:"product_id,type:uuid,notnull"`
	Product            Product         `bun:"product,rel:has-one,join:product_id=id"`
	ProductVariationID *uuid.UUID      `bun:"product_variation_id,type:uuid"`
	CurrentStock       decimal.Decimal `bun:"current_stock,type:decimal(10,3),notnull"`
	MinStock           decimal.Decimal `bun:"min_stock,type:decimal(10,3),notnull"`
	MaxStock           decimal.Decimal `bun:"max_stock,type:decimal(10,3),notnull"`
	Unit               string          `bun:"unit,notnull"`
	IsActive           bool            `bun:"is_active,notnull"`
}

// FromDomain converte domain para model
func (s *Stock) FromDomain(stock *stockentity.Stock) {
	if stock == nil {
		return
	}
	*s = Stock{
		Entity: entitymodel.FromDomain(stock.Entity),
		StockCommonAttributes: StockCommonAttributes{
			ProductID:          stock.ProductID,
			ProductVariationID: stock.ProductVariationID,
			CurrentStock:       stock.CurrentStock,
			MinStock:           stock.MinStock,
			MaxStock:           stock.MaxStock,
			Unit:               stock.Unit,
			IsActive:           stock.IsActive,
		},
	}
}

// ToDomain converte model para domain
func (s *Stock) ToDomain() *stockentity.Stock {
	if s == nil {
		return nil
	}
	return &stockentity.Stock{
		Entity: s.Entity.ToDomain(),
		StockCommonAttributes: stockentity.StockCommonAttributes{
			ProductID:          s.ProductID,
			Product:            *s.Product.ToDomain(),
			ProductVariationID: s.ProductVariationID,
			CurrentStock:       s.CurrentStock,
			MinStock:           s.MinStock,
			MaxStock:           s.MaxStock,
			Unit:               s.Unit,
			IsActive:           s.IsActive,
		},
	}
}
