package stockdto

import (
	"time"

	"github.com/google/uuid"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockAlertDTO representa o DTO de alerta de estoque
type StockAlertDTO struct {
	ID         uuid.UUID  `json:"id"`
	StockID    uuid.UUID  `json:"stock_id"`
	Type       string     `json:"type"`
	Message    string     `json:"message"`
	IsResolved bool       `json:"is_resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *uuid.UUID `json:"resolved_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// FromDomain converte domain para DTO
func (sa *StockAlertDTO) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	*sa = StockAlertDTO{
		ID:         alert.ID,
		StockID:    alert.StockID,
		Type:       string(alert.Type),
		Message:    alert.Message,
		IsResolved: alert.IsResolved,
		ResolvedAt: alert.ResolvedAt,
		ResolvedBy: alert.ResolvedBy,
		CreatedAt:  alert.CreatedAt,
	}
}
