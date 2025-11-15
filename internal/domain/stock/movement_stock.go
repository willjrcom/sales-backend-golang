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
	StockID    uuid.UUID
	Type       MovementType
	Quantity   decimal.Decimal
	Reason     string
	OrderID    *uuid.UUID
	EmployeeID uuid.UUID
	Price      decimal.Decimal
	TotalPrice decimal.Decimal
}

// MovementType define o tipo de movimento de estoque
type MovementType string

const (
	MovementTypeIn        MovementType = "in"         // Entrada de estoque
	MovementTypeOut       MovementType = "out"        // Saída de estoque
	MovementTypeAdjustIn  MovementType = "adjust_in"  // Ajuste manual
	MovementTypeAdjustOut MovementType = "adjust_out" // Ajuste manual

)
