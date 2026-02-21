package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// StockAlert model
type StockAlert struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_alerts,alias:stock_alert"`
	StockAlertCommonAttributes
}

type AlertType string

type StockAlertCommonAttributes struct {
	StockID            uuid.UUID         `bun:"stock_id,type:uuid,notnull"`
	Type               AlertType         `bun:"type,notnull"`
	Message            string            `bun:"message,notnull"`
	IsResolved         bool              `bun:"is_resolved,notnull"`
	ResolvedAt         *time.Time        `bun:"resolved_at"`
	ResolvedBy         *uuid.UUID        `bun:"resolved_by,type:uuid"`
	ProductID          uuid.UUID         `bun:"product_id,type:uuid,notnull"`
	Product            *Product          `bun:"rel:belongs-to"`
	ProductVariationID uuid.UUID         `bun:"product_variation_id,type:uuid,notnull"`
	ProductVariation   *ProductVariation `bun:"rel:belongs-to"`
}

// FromDomain converte domain para model
func (s *StockAlert) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	*s = StockAlert{
		Entity: entitymodel.FromDomain(alert.Entity),
		StockAlertCommonAttributes: StockAlertCommonAttributes{
			StockID:            alert.StockID,
			Type:               AlertType(alert.Type),
			Message:            alert.Message,
			IsResolved:         alert.IsResolved,
			ResolvedAt:         alert.ResolvedAt,
			ResolvedBy:         alert.ResolvedBy,
			ProductID:          alert.ProductID,
			ProductVariationID: alert.ProductVariationID,
		},
	}
}

// ToDomain converte model para domain
func (s *StockAlert) ToDomain() *stockentity.StockAlert {
	if s == nil {
		return nil
	}

	var product *productentity.Product
	if s.Product != nil {
		product = s.Product.ToDomain()
	}

	var productVariation *productentity.ProductVariation
	if s.ProductVariation != nil {
		v := s.ProductVariation.ToDomain()
		productVariation = &v
	}

	return &stockentity.StockAlert{
		Entity: s.Entity.ToDomain(),
		StockAlertCommonAttributes: stockentity.StockAlertCommonAttributes{
			StockID:            s.StockID,
			Type:               stockentity.AlertType(s.Type),
			Message:            s.Message,
			IsResolved:         s.IsResolved,
			ResolvedAt:         s.ResolvedAt,
			ResolvedBy:         s.ResolvedBy,
			ProductID:          s.ProductID,
			Product:            product,
			ProductVariationID: s.ProductVariationID,
			ProductVariation:   productVariation,
		},
	}
}
