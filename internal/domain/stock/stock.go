package stockentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrStockMovementRequired = errors.New("stock movement is required")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrProductNotFound       = errors.New("product not found")
)

// Stock representa o estoque atual de um produto
type Stock struct {
	entity.Entity
	StockCommonAttributes
}

type StockCommonAttributes struct {
	ProductID    uuid.UUID       `json:"product_id"`
	CurrentStock decimal.Decimal `json:"current_stock"`
	MinStock     decimal.Decimal `json:"min_stock"` // Estoque mínimo para alertas
	MaxStock     decimal.Decimal `json:"max_stock"` // Estoque máximo
	Unit         string          `json:"unit"`      // Unidade (kg, unidades, etc)
	IsActive     bool            `json:"is_active"` // Se o controle de estoque está ativo
}

// StockMovement representa um movimento de estoque (entrada/saída)
type StockMovement struct {
	entity.Entity
	StockMovementCommonAttributes
}

type StockMovementCommonAttributes struct {
	StockID      uuid.UUID       `json:"stock_id"`
	ProductID    uuid.UUID       `json:"product_id"`
	Type         MovementType    `json:"type"`
	Quantity     decimal.Decimal `json:"quantity"`
	Reason       string          `json:"reason"`
	OrderID      *uuid.UUID      `json:"order_id,omitempty"`      // Se relacionado a um pedido
	OrderNumber  *int            `json:"order_number,omitempty"`  // Número do pedido
	EmployeeID   *uuid.UUID      `json:"employee_id,omitempty"`   // Funcionário responsável
	EmployeeName *string         `json:"employee_name,omitempty"` // Nome do funcionário
	UnitCost     decimal.Decimal `json:"unit_cost"`               // Custo unitário no momento do movimento
	TotalCost    decimal.Decimal `json:"total_cost"`              // Custo total do movimento
	Notes        string          `json:"notes"`                   // Observações adicionais
}

// MovementType define o tipo de movimento de estoque
type MovementType string

const (
	MovementTypeIn     MovementType = "in"     // Entrada de estoque
	MovementTypeOut    MovementType = "out"    // Saída de estoque
	MovementTypeAdjust MovementType = "adjust" // Ajuste manual
)

// StockAlert representa alertas de estoque (baixo estoque, etc)
type StockAlert struct {
	entity.Entity
	StockAlertCommonAttributes
}

type StockAlertCommonAttributes struct {
	StockID    uuid.UUID  `json:"stock_id"`
	ProductID  uuid.UUID  `json:"product_id"`
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

// NewStock cria um novo estoque
func NewStock(productID uuid.UUID, initialStock, minStock, maxStock decimal.Decimal, unit string) *Stock {
	return &Stock{
		Entity: entity.NewEntity(),
		StockCommonAttributes: StockCommonAttributes{
			ProductID:    productID,
			CurrentStock: initialStock,
			MinStock:     minStock,
			MaxStock:     maxStock,
			Unit:         unit,
			IsActive:     true,
		},
	}
}

// AddStock adiciona estoque manualmente
func (s *Stock) AddStock(quantity decimal.Decimal, reason string, employeeID *uuid.UUID, unitCost decimal.Decimal, notes string) (*StockMovement, error) {
	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	s.CurrentStock = s.CurrentStock.Add(quantity)

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    s.ID,
			ProductID:  s.ProductID,
			Type:       MovementTypeIn,
			Quantity:   quantity,
			Reason:     reason,
			EmployeeID: employeeID,
			UnitCost:   unitCost,
			TotalCost:  unitCost.Mul(quantity),
			Notes:      notes,
		},
	}

	return movement, nil
}

// RemoveStock remove estoque manualmente
func (s *Stock) RemoveStock(quantity decimal.Decimal, reason string, employeeID *uuid.UUID, unitCost decimal.Decimal, notes string) (*StockMovement, error) {
	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	if s.CurrentStock.LessThan(quantity) {
		return nil, ErrInsufficientStock
	}

	s.CurrentStock = s.CurrentStock.Sub(quantity)

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    s.ID,
			ProductID:  s.ProductID,
			Type:       MovementTypeOut,
			Quantity:   quantity,
			Reason:     reason,
			EmployeeID: employeeID,
			UnitCost:   unitCost,
			TotalCost:  unitCost.Mul(quantity),
			Notes:      notes,
		},
	}

	return movement, nil
}

// AdjustStock ajusta o estoque para um valor específico
func (s *Stock) AdjustStock(newStock decimal.Decimal, reason string, employeeID *uuid.UUID, unitCost decimal.Decimal, notes string) (*StockMovement, error) {
	if newStock.LessThan(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	difference := newStock.Sub(s.CurrentStock)
	s.CurrentStock = newStock

	movementType := MovementTypeAdjust
	if difference.GreaterThan(decimal.Zero) {
		movementType = MovementTypeIn
	} else if difference.LessThan(decimal.Zero) {
		movementType = MovementTypeOut
		difference = difference.Abs()
	} else {
		// Não há diferença, não criar movimento
		return nil, nil
	}

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    s.ID,
			ProductID:  s.ProductID,
			Type:       movementType,
			Quantity:   difference,
			Reason:     reason,
			EmployeeID: employeeID,
			UnitCost:   unitCost,
			TotalCost:  unitCost.Mul(difference),
			Notes:      notes,
		},
	}

	return movement, nil
}

// ReserveStock reserva estoque para um pedido (quando pedido fica pending)
func (s *Stock) ReserveStock(quantity decimal.Decimal, orderID uuid.UUID, orderNumber int, unitCost decimal.Decimal) (*StockMovement, error) {
	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	if s.CurrentStock.LessThan(quantity) {
		return nil, ErrInsufficientStock
	}

	s.CurrentStock = s.CurrentStock.Sub(quantity)

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:     s.ID,
			ProductID:   s.ProductID,
			Type:        MovementTypeOut,
			Quantity:    quantity,
			Reason:      "Reserva automática - Pedido pendente",
			OrderID:     &orderID,
			OrderNumber: &orderNumber,
			UnitCost:    unitCost,
			TotalCost:   unitCost.Mul(quantity),
			Notes:       "Débito automático ao deixar pedido pendente",
		},
	}

	return movement, nil
}

// RestoreStock restaura estoque quando pedido é cancelado
func (s *Stock) RestoreStock(quantity decimal.Decimal, orderID uuid.UUID, orderNumber int, unitCost decimal.Decimal) (*StockMovement, error) {
	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	s.CurrentStock = s.CurrentStock.Add(quantity)

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:     s.ID,
			ProductID:   s.ProductID,
			Type:        MovementTypeIn,
			Quantity:    quantity,
			Reason:      "Restauração automática - Pedido cancelado",
			OrderID:     &orderID,
			OrderNumber: &orderNumber,
			UnitCost:    unitCost,
			TotalCost:   unitCost.Mul(quantity),
			Notes:       "Crédito automático ao cancelar pedido",
		},
	}

	return movement, nil
}

// CheckAlerts verifica se há alertas de estoque
func (s *Stock) CheckAlerts() []*StockAlert {
	var alerts []*StockAlert

	if !s.IsActive {
		return alerts
	}

	// Alerta de estoque baixo
	if s.CurrentStock.LessThanOrEqual(s.MinStock) && s.CurrentStock.GreaterThan(decimal.Zero) {
		alerts = append(alerts, &StockAlert{
			Entity: entity.NewEntity(),
			StockAlertCommonAttributes: StockAlertCommonAttributes{
				StockID:   s.ID,
				ProductID: s.ProductID,

				Type:       AlertTypeLowStock,
				Message:    "Estoque abaixo do mínimo",
				IsResolved: false,
			},
		})
	}

	// Alerta de estoque zerado
	if s.CurrentStock.LessThanOrEqual(decimal.Zero) {
		alerts = append(alerts, &StockAlert{
			Entity: entity.NewEntity(),
			StockAlertCommonAttributes: StockAlertCommonAttributes{
				StockID:   s.ID,
				ProductID: s.ProductID,

				Type:       AlertTypeOutOfStock,
				Message:    "Produto sem estoque",
				IsResolved: false,
			},
		})
	}

	// Alerta de estoque acima do máximo
	if s.MaxStock.GreaterThan(decimal.Zero) && s.CurrentStock.GreaterThan(s.MaxStock) {
		alerts = append(alerts, &StockAlert{
			Entity: entity.NewEntity(),
			StockAlertCommonAttributes: StockAlertCommonAttributes{
				StockID:   s.ID,
				ProductID: s.ProductID,

				Type:       AlertTypeOverStock,
				Message:    "Estoque acima do máximo",
				IsResolved: false,
			},
		})
	}

	return alerts
}

// IsLowStock verifica se o estoque está baixo
func (s *Stock) IsLowStock() bool {
	return s.CurrentStock.LessThanOrEqual(s.MinStock)
}

// IsOutOfStock verifica se o produto está sem estoque
func (s *Stock) IsOutOfStock() bool {
	return s.CurrentStock.LessThanOrEqual(decimal.Zero)
}

// GetStockLevel retorna o nível de estoque em porcentagem
func (s *Stock) GetStockLevel() decimal.Decimal {
	if s.MaxStock.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	return s.CurrentStock.Div(s.MaxStock).Mul(decimal.NewFromInt(100))
}
