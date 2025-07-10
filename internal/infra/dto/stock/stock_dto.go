package stockdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

// StockDTO representa o DTO de estoque
type StockDTO struct {
	ID           uuid.UUID       `json:"id"`
	ProductID    uuid.UUID       `json:"product_id"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"`
	MaxStock     decimal.Decimal `json:"max_stock"`
	Unit         string          `json:"unit"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// StockCreateDTO representa o DTO para criar estoque
type StockCreateDTO struct {
	ProductID    uuid.UUID       `json:"product_id"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"`
	MaxStock     decimal.Decimal `json:"max_stock"`
	Unit         string          `json:"unit"`
	IsActive     bool            `json:"is_active"`
}

// StockUpdateDTO representa o DTO para atualizar estoque
type StockUpdateDTO struct {
	CurrentStock *decimal.Decimal `json:"current_stock,omitempty"`
	MinStock     *decimal.Decimal `json:"min_stock,omitempty"`
	MaxStock     *decimal.Decimal `json:"max_stock,omitempty"`
	Unit         *string          `json:"unit,omitempty"`
	IsActive     *bool            `json:"is_active,omitempty"`
}

// StockMovementDTO representa o DTO de movimento de estoque
type StockMovementDTO struct {
	ID           uuid.UUID       `json:"id"`
	StockID      uuid.UUID       `json:"stock_id"`
	ProductID    uuid.UUID       `json:"product_id"`
	Type         string          `json:"type"`
	Quantity     decimal.Decimal `json:"quantity"`
	Reason       string          `json:"reason"`
	OrderID      *uuid.UUID      `json:"order_id,omitempty"`
	OrderNumber  *int            `json:"order_number,omitempty"`
	EmployeeID   *uuid.UUID      `json:"employee_id,omitempty"`
	EmployeeName *string         `json:"employee_name,omitempty"`
	UnitCost     decimal.Decimal `json:"unit_cost"`
	TotalCost    decimal.Decimal `json:"total_cost"`
	Notes        string          `json:"notes"`
	CreatedAt    time.Time       `json:"created_at"`
}

// StockMovementCreateDTO representa o DTO para criar movimento
type StockMovementCreateDTO struct {
	StockID      uuid.UUID       `json:"stock_id"`
	ProductID    uuid.UUID       `json:"product_id"`
	Type         string          `json:"type"`
	Quantity     decimal.Decimal `json:"quantity"`
	Reason       string          `json:"reason"`
	OrderID      *uuid.UUID      `json:"order_id,omitempty"`
	OrderNumber  *int            `json:"order_number,omitempty"`
	EmployeeID   *uuid.UUID      `json:"employee_id,omitempty"`
	EmployeeName *string         `json:"employee_name,omitempty"`
	UnitCost     decimal.Decimal `json:"unit_cost"`
	Notes        string          `json:"notes"`
}

// StockAlertDTO representa o DTO de alerta de estoque
type StockAlertDTO struct {
	ID         uuid.UUID  `json:"id"`
	StockID    uuid.UUID  `json:"stock_id"`
	ProductID  uuid.UUID  `json:"product_id"`
	Type       string     `json:"type"`
	Message    string     `json:"message"`
	IsResolved bool       `json:"is_resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *uuid.UUID `json:"resolved_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// StockReportDTO representa o DTO para relatórios de estoque
type StockReportDTO struct {
	ProductID    uuid.UUID       `json:"product_id"`
	ProductName  string          `json:"product_name"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"`
	MaxStock     decimal.Decimal `json:"max_stock"`
	Unit         string          `json:"unit"`
	StockLevel   decimal.Decimal `json:"stock_level"` // Porcentagem
	IsLowStock   bool            `json:"is_low_stock"`
	IsOutOfStock bool            `json:"is_out_of_stock"`
	TotalIn      decimal.Decimal `json:"total_in"`
	TotalOut     decimal.Decimal `json:"total_out"`
	TotalCost    decimal.Decimal `json:"total_cost"`
}

// StockWithProductDTO representa o DTO de estoque com informações do produto
type StockWithProductDTO struct {
	StockDTO
	ProductName string `json:"product_name,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
}

// StockMovementReportDTO representa o DTO para relatório de movimentos
type StockMovementReportDTO struct {
	Date         time.Time       `json:"date"`
	ProductID    uuid.UUID       `json:"product_id"`
	ProductName  string          `json:"product_name"`
	Type         string          `json:"type"`
	Quantity     decimal.Decimal `json:"quantity"`
	Reason       string          `json:"reason"`
	OrderNumber  *int            `json:"order_number,omitempty"`
	EmployeeName *string         `json:"employee_name,omitempty"`
	UnitCost     decimal.Decimal `json:"unit_cost"`
	TotalCost    decimal.Decimal `json:"total_cost"`
}

// StockReportSummaryDTO representa o resumo do relatório de estoque
type StockReportSummaryDTO struct {
	TotalProducts     int             `json:"total_products"`
	TotalLowStock     int             `json:"total_low_stock"`
	TotalOutOfStock   int             `json:"total_out_of_stock"`
	TotalActiveAlerts int             `json:"total_active_alerts"`
	TotalStockValue   decimal.Decimal `json:"total_stock_value"`
}

// StockReportCompleteDTO representa o relatório completo de estoque
type StockReportCompleteDTO struct {
	Summary            StockReportSummaryDTO `json:"summary"`
	AllStocks          []StockDTO            `json:"all_stocks"`
	LowStockProducts   []StockDTO            `json:"low_stock_products"`
	OutOfStockProducts []StockDTO            `json:"out_of_stock_products"`
	ActiveAlerts       []StockAlertDTO       `json:"active_alerts"`
	GeneratedAt        time.Time             `json:"generated_at"`
}

// FromDomain converte domain para DTO
func (s *StockDTO) FromDomain(stock *stockentity.Stock) {
	if stock == nil {
		return
	}
	*s = StockDTO{
		ID:           stock.ID,
		ProductID:    stock.ProductID,
		CurrentStock: stock.CurrentStock,
		MinStock:     stock.MinStock,
		MaxStock:     stock.MaxStock,
		Unit:         stock.Unit,
		IsActive:     stock.IsActive,
		CreatedAt:    stock.CreatedAt,
		UpdatedAt:    stock.UpdatedAt,
	}
}

// ToDomain converte DTO para domain
func (s *StockCreateDTO) ToDomain() *stockentity.Stock {
	return stockentity.NewStock(
		s.ProductID,
		s.CurrentStock,
		s.MinStock,
		s.MaxStock,
		s.Unit,
	)
}

// FromDomain converte domain para DTO
func (sm *StockMovementDTO) FromDomain(movement *stockentity.StockMovement) {
	if movement == nil {
		return
	}
	*sm = StockMovementDTO{
		ID:           movement.ID,
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
		CreatedAt:    movement.CreatedAt,
	}
}

// ToDomain converte DTO para domain
func (sm *StockMovementCreateDTO) ToDomain() *stockentity.StockMovement {
	return &stockentity.StockMovement{
		Entity: entity.NewEntity(),
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
			TotalCost:    sm.UnitCost.Mul(sm.Quantity),
			Notes:        sm.Notes,
		},
	}
}

// FromDomain converte domain para DTO
func (sa *StockAlertDTO) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	*sa = StockAlertDTO{
		ID:         alert.ID,
		StockID:    alert.StockID,
		ProductID:  alert.ProductID,
		Type:       string(alert.Type),
		Message:    alert.Message,
		IsResolved: alert.IsResolved,
		ResolvedAt: alert.ResolvedAt,
		ResolvedBy: alert.ResolvedBy,
		CreatedAt:  alert.CreatedAt,
	}
}
