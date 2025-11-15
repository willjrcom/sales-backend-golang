package stockentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// StockAlert representa alertas de estoque (baixo estoque, etc)
type StockAlert struct {
	entity.Entity
	StockAlertCommonAttributes
}

type StockAlertCommonAttributes struct {
	StockID    uuid.UUID  `json:"stock_id"`
	Type       AlertType  `json:"type"`
	Message    string     `json:"message"`
	IsResolved bool       `json:"is_resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *uuid.UUID `json:"resolved_by,omitempty"`
}

// AlertType define o tipo de alerta
type AlertType string

const (
	AlertTypeLowStock   AlertType = "low_stock"    // Estoque baixo
	AlertTypeOutOfStock AlertType = "out_of_stock" // Sem estoque
	AlertTypeOverStock  AlertType = "over_stock"   // Estoque acima do máximo
)
