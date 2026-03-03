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
	BatchID    *uuid.UUID // ID do lote associado (opcional)
	Type       MovementType
	Quantity   decimal.Decimal
	Reason     string
	OrderID    *uuid.UUID
	EmployeeID uuid.UUID
	Price      decimal.Decimal // Entrada: Custo do lote | Saída: Preço de venda
}

// MovementType define o tipo de movimento de estoque
type MovementType string

const (
	MovementTypeIn        MovementType = "in"         // Entrada de estoque
	MovementTypeOut       MovementType = "out"        // Saída de estoque
	MovementTypeAdjustIn  MovementType = "adjust_in"  // Ajuste manual
	MovementTypeAdjustOut MovementType = "adjust_out" // Ajuste manual
	MovementTypeReserve   MovementType = "reserve"    // Reserva para pedido
	MovementTypeRestore   MovementType = "restore"    // Restauração de reserva
)

func NewStockMovement(stockID uuid.UUID, quantity decimal.Decimal, reason string, employeeID uuid.UUID, price decimal.Decimal) *StockMovement {
	return &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    stockID,
			Quantity:   quantity,
			Reason:     reason,
			EmployeeID: employeeID,
			Price:      price,
		},
	}
}
