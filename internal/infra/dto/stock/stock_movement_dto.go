package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockMovementDTO representa o DTO de movimento de estoque
type StockMovementDTO struct {
	ID         uuid.UUID       `json:"id"`
	StockID    uuid.UUID       `json:"stock_id"`
	Type       string          `json:"type"`
	Reason     string          `json:"reason"`
	OrderID    *uuid.UUID      `json:"order_id,omitempty"`
	EmployeeID uuid.UUID       `json:"employee_id,omitempty"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"unit_cost"`
	TotalPrice decimal.Decimal `json:"total_cost"`
	CreatedAt  time.Time       `json:"created_at"`
}

// FromDomain converte domain para DTO
func (sm *StockMovementDTO) FromDomain(movement *stockentity.StockMovement) {
	if movement == nil {
		return
	}
	*sm = StockMovementDTO{
		ID:         movement.ID,
		StockID:    movement.StockID,
		Type:       string(movement.Type),
		Quantity:   movement.Quantity,
		Reason:     movement.Reason,
		OrderID:    movement.OrderID,
		EmployeeID: movement.EmployeeID,
		Price:      movement.Price,
		TotalPrice: movement.TotalPrice,
		CreatedAt:  movement.CreatedAt,
	}
}
