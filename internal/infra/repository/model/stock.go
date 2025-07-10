package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// Stock model
type Stock struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stocks,alias:stock"`
	StockCommonAttributes
}

type StockCommonAttributes struct {
	ProductID    uuid.UUID       `bun:"product_id,type:uuid,notnull"`
	CurrentStock decimal.Decimal `bun:"current_stock,type:decimal(10,3),notnull"`
	MinStock     decimal.Decimal `bun:"min_stock,type:decimal(10,3),notnull"`
	MaxStock     decimal.Decimal `bun:"max_stock,type:decimal(10,3),notnull"`
	Unit         string          `bun:"unit,notnull"`
	IsActive     bool            `bun:"is_active,notnull"`
}

// StockMovement model
type StockMovement struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_movements,alias:stock_movement"`
	StockMovementCommonAttributes
}

type StockMovementCommonAttributes struct {
	StockID      uuid.UUID       `bun:"stock_id,type:uuid,notnull"`
	ProductID    uuid.UUID       `bun:"product_id,type:uuid,notnull"`
	Type         string          `bun:"type,notnull"`
	Quantity     decimal.Decimal `bun:"quantity,type:decimal(10,3),notnull"`
	Reason       string          `bun:"reason,notnull"`
	OrderID      *uuid.UUID      `bun:"order_id,type:uuid"`
	OrderNumber  *int            `bun:"order_number"`
	EmployeeID   *uuid.UUID      `bun:"employee_id,type:uuid"`
	EmployeeName *string         `bun:"employee_name"`
	UnitCost     decimal.Decimal `bun:"unit_cost,type:decimal(10,2),notnull"`
	TotalCost    decimal.Decimal `bun:"total_cost,type:decimal(10,2),notnull"`
	Notes        string          `bun:"notes"`
}

// StockAlert model
type StockAlert struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_alerts,alias:stock_alert"`
	StockAlertCommonAttributes
}

type StockAlertCommonAttributes struct {
	StockID    uuid.UUID  `bun:"stock_id,type:uuid,notnull"`
	ProductID  uuid.UUID  `bun:"product_id,type:uuid,notnull"`
	Type       string     `bun:"type,notnull"`
	Message    string     `bun:"message,notnull"`
	IsResolved bool       `bun:"is_resolved,notnull"`
	ResolvedAt *time.Time `bun:"resolved_at"`
	ResolvedBy *uuid.UUID `bun:"resolved_by,type:uuid"`
}

// FromDomain converte domain para model
func (s *Stock) FromDomain(stock *stockentity.Stock) {
	if stock == nil {
		return
	}
	*s = Stock{
		Entity: entitymodel.FromDomain(stock.Entity),
		StockCommonAttributes: StockCommonAttributes{
			ProductID:    stock.ProductID,
			CurrentStock: stock.CurrentStock,
			MinStock:     stock.MinStock,
			MaxStock:     stock.MaxStock,
			Unit:         stock.Unit,
			IsActive:     stock.IsActive,
		},
	}
}

// ToDomain converte model para domain
func (s *Stock) ToDomain() *stockentity.Stock {
	if s == nil {
		return nil
	}
	return &stockentity.Stock{
		Entity: s.Entity.ToDomain(),
		StockCommonAttributes: stockentity.StockCommonAttributes{
			ProductID:    s.ProductID,
			CurrentStock: s.CurrentStock,
			MinStock:     s.MinStock,
			MaxStock:     s.MaxStock,
			Unit:         s.Unit,
			IsActive:     s.IsActive,
		},
	}
}

// FromDomain converte domain para model
func (sm *StockMovement) FromDomain(movement *stockentity.StockMovement) {
	if movement == nil {
		return
	}
	*sm = StockMovement{
		Entity: entitymodel.FromDomain(movement.Entity),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:      movement.StockID,
			ProductID:    movement.ProductID,
			Type:         string(movement.Type),
			Quantity:     movement.Quantity,
			Reason:       movement.Reason,
			OrderID:      movement.OrderID,
			OrderNumber:  movement.OrderNumber,
			EmployeeID:   movement.EmployeeID,
			EmployeeName: movement.EmployeeName,
			UnitCost:     movement.UnitCost,
			TotalCost:    movement.TotalCost,
			Notes:        movement.Notes,
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
			StockID:      sm.StockID,
			ProductID:    sm.ProductID,
			Type:         stockentity.MovementType(sm.Type),
			Quantity:     sm.Quantity,
			Reason:       sm.Reason,
			OrderID:      sm.OrderID,
			OrderNumber:  sm.OrderNumber,
			EmployeeID:   sm.EmployeeID,
			EmployeeName: sm.EmployeeName,
			UnitCost:     sm.UnitCost,
			TotalCost:    sm.TotalCost,
			Notes:        sm.Notes,
		},
	}
}

// FromDomain converte domain para model
func (sa *StockAlert) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	*sa = StockAlert{
		Entity: entitymodel.FromDomain(alert.Entity),
		StockAlertCommonAttributes: StockAlertCommonAttributes{
			StockID:    alert.StockID,
			ProductID:  alert.ProductID,
			Type:       string(alert.Type),
			Message:    alert.Message,
			IsResolved: alert.IsResolved,
			ResolvedAt: alert.ResolvedAt,
			ResolvedBy: alert.ResolvedBy,
		},
	}
}

// ToDomain converte model para domain
func (sa *StockAlert) ToDomain() *stockentity.StockAlert {
	if sa == nil {
		return nil
	}
	return &stockentity.StockAlert{
		Entity: sa.Entity.ToDomain(),
		StockAlertCommonAttributes: stockentity.StockAlertCommonAttributes{
			StockID:    sa.StockID,
			ProductID:  sa.ProductID,
			Type:       stockentity.AlertType(sa.Type),
			Message:    sa.Message,
			IsResolved: sa.IsResolved,
			ResolvedAt: sa.ResolvedAt,
			ResolvedBy: sa.ResolvedBy,
		},
	}
}

// ToDTO converte model para DTO
func (s *Stock) ToDTO() *stockdto.StockDTO {
	if s == nil {
		return nil
	}
	return &stockdto.StockDTO{
		ID:           s.ID,
		ProductID:    s.ProductID,
		CurrentStock: s.CurrentStock,
		MinStock:     s.MinStock,
		MaxStock:     s.MaxStock,
		Unit:         s.Unit,
		IsActive:     s.IsActive,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

// ToDTO converte model para DTO
func (sa *StockAlert) ToDTO() *stockdto.StockAlertDTO {
	if sa == nil {
		return nil
	}
	return &stockdto.StockAlertDTO{
		ID:         sa.ID,
		StockID:    sa.StockID,
		ProductID:  sa.ProductID,
		Type:       sa.Type,
		Message:    sa.Message,
		IsResolved: sa.IsResolved,
		ResolvedAt: sa.ResolvedAt,
		ResolvedBy: sa.ResolvedBy,
		CreatedAt:  sa.CreatedAt,
	}
}
