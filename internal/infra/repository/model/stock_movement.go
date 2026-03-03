package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// StockMovement model
type StockMovement struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_movements,alias:stock_movement"`
	StockMovementCommonAttributes
}

type StockMovementCommonAttributes struct {
	StockID    uuid.UUID        `bun:"stock_id,type:uuid,notnull"`
	BatchID    *uuid.UUID       `bun:"batch_id,type:uuid"`
	Type       string           `bun:"type,notnull"`
	Quantity   *decimal.Decimal `bun:"quantity,type:decimal(10,3),notnull"`
	Reason     string           `bun:"reason,notnull"`
	OrderID    *uuid.UUID       `bun:"order_id"`
	EmployeeID uuid.UUID        `bun:"employee_id,notnull"`
	Price      *decimal.Decimal `bun:"price,type:decimal(10,2),notnull"`
}

// FromDomain converte domain para model
func (sm *StockMovement) FromDomain(movement *stockentity.StockMovement) {
	if movement == nil {
		return
	}
	sm.ID = movement.ID
	sm.StockID = movement.StockID
	sm.BatchID = movement.BatchID
	sm.Type = string(movement.Type)
	sm.Quantity = &movement.Quantity
	sm.Reason = movement.Reason
	sm.OrderID = movement.OrderID
	sm.EmployeeID = movement.EmployeeID
	sm.Price = &movement.Price
}

// ToDomain converte model para domain
func (sm *StockMovement) ToDomain() *stockentity.StockMovement {
	if sm == nil {
		return nil
	}
	return &stockentity.StockMovement{
		Entity: sm.Entity.ToDomain(),
		StockMovementCommonAttributes: stockentity.StockMovementCommonAttributes{
			StockID:    sm.StockID,
			BatchID:    sm.BatchID,
			Type:       stockentity.MovementType(sm.Type),
			Quantity:   sm.GetQuantity(),
			Reason:     sm.Reason,
			OrderID:    sm.OrderID,
			EmployeeID: sm.EmployeeID,
			Price:      sm.GetPrice(),
		},
	}
}

func (sm *StockMovement) GetQuantity() decimal.Decimal {
	if sm.Quantity == nil {
		return decimal.Zero
	}
	return *sm.Quantity
}

func (sm *StockMovement) GetPrice() decimal.Decimal {
	if sm.Price == nil {
		return decimal.Zero
	}
	return *sm.Price
}
