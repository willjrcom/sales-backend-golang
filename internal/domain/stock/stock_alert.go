package stockentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

// StockAlert representa alertas de estoque (baixo estoque, etc)
type StockAlert struct {
	entity.Entity
	StockAlertCommonAttributes
}

type StockAlertCommonAttributes struct {
	StockID            uuid.UUID
	Type               AlertType
	Message            string
	IsResolved         bool
	ResolvedAt         *time.Time
	ResolvedBy         *uuid.UUID
	ProductID          uuid.UUID
	Product            *productentity.Product
	ProductVariationID uuid.UUID
	ProductVariation   *productentity.ProductVariation
}

// AlertType define o tipo de alerta
type AlertType string

const (
	AlertTypeLowStock   AlertType = "low_stock"    // Estoque baixo
	AlertTypeOutOfStock AlertType = "out_of_stock" // Sem estoque
	AlertTypeOverStock  AlertType = "over_stock"   // Estoque acima do máximo
)
