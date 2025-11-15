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
	StockID    uuid.UUID       `bun:"stock_id,type:uuid,notnull"`
	Type       string          `bun:"type,notnull"`
	Quantity   decimal.Decimal `bun:"quantity,type:decimal(10,3),notnull"`
	Reason     string          `bun:"reason,notnull"`
	OrderID    *uuid.UUID      `bun:"order_id,type:uuid"`
	EmployeeID uuid.UUID       `bun:"employee_id,type:uuid"`
	Price      decimal.Decimal `bun:"unit_cost,type:decimal(10,2),notnull"`
	TotalPrice decimal.Decimal `bun:"total_cost,type:decimal(10,2),notnull"`
}

// FromDomain converte domain para model
func (sm *StockMovement) FromDomain(movement *stockentity.StockMovement) {
	if movement == nil {
		return
	}
	*sm = StockMovement{
		Entity: entitymodel.FromDomain(movement.Entity),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    movement.StockID,
			Type:       string(movement.Type),
			Quantity:   movement.Quantity,
			Reason:     movement.Reason,
			OrderID:    movement.OrderID,
			EmployeeID: movement.EmployeeID,
			Price:      movement.Price,
			TotalPrice: movement.TotalPrice,
		},
	}
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
			Type:       stockentity.MovementType(sm.Type),
			Quantity:   sm.Quantity,
			Reason:     sm.Reason,
			OrderID:    sm.OrderID,
			EmployeeID: sm.EmployeeID,
			Price:      sm.Price,
			TotalPrice: sm.TotalPrice,
		},
	}
}
