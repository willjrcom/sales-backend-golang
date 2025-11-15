package stockentity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

// StockMovement representa um movimento de estoque (entrada/saída)
type StockMovement struct {
	entity.Entity
	StockMovementCommonAttributes
}

type StockMovementCommonAttributes struct {
	StockID    uuid.UUID       `json:"stock_id"`
	Type       MovementType    `json:"type"`
	Quantity   decimal.Decimal `json:"quantity"`
	Reason     string          `json:"reason"`
	OrderID    *uuid.UUID      `json:"order_id,omitempty"`    // Se relacionado a um pedido
	EmployeeID uuid.UUID       `json:"employee_id,omitempty"` // Funcionário responsável
	Price      decimal.Decimal `json:"price"`                 // Custo unitário no momento do movimento
	TotalPrice decimal.Decimal `json:"total_price"`           // Custo total do movimento
}

// MovementType define o tipo de movimento de estoque
type MovementType string

const (
	MovementTypeIn        MovementType = "in"         // Entrada de estoque
	MovementTypeOut       MovementType = "out"        // Saída de estoque
	MovementTypeAdjustIn  MovementType = "adjust_in"  // Ajuste manual
	MovementTypeAdjustOut MovementType = "adjust_out" // Ajuste manual

)
