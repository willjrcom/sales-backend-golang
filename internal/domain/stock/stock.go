package stockentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
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
	ProductID    uuid.UUID
	Product      productentity.Product
	CurrentStock decimal.Decimal
	MinStock     decimal.Decimal
	MaxStock     decimal.Decimal
	Unit         string
	IsActive     bool
}

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

// AddMovementStock adiciona estoque manualmente
func (s *Stock) AddMovementStock(quantity decimal.Decimal, reason string, employeeID uuid.UUID, price decimal.Decimal, totalPrice decimal.Decimal) (*StockMovement, error) {
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
			Type:       MovementTypeIn,
			Quantity:   quantity,
			Reason:     reason,
			EmployeeID: employeeID,
			Price:      price,
			TotalPrice: totalPrice,
		},
	}

	return movement, nil
}

// RemoveMovementStock remove estoque manualmente
func (s *Stock) RemoveMovementStock(quantity decimal.Decimal, reason string, employeeID uuid.UUID, price decimal.Decimal, totalPrice decimal.Decimal) (*StockMovement, error) {
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
			Type:       MovementTypeOut,
			Quantity:   quantity,
			Reason:     reason,
			EmployeeID: employeeID,
			Price:      price,
			TotalPrice: totalPrice,
		},
	}

	return movement, nil
}

// AdjustMovementStock ajusta o estoque para um valor específico
func (s *Stock) AdjustMovementStock(newStock decimal.Decimal, reason string, employeeID uuid.UUID) (*StockMovement, error) {
	if newStock.LessThan(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if !s.IsActive {
		return nil, errors.New("stock control is not active")
	}

	difference := newStock.Sub(s.CurrentStock)
	s.CurrentStock = newStock

	var movementType MovementType
	if difference.GreaterThan(decimal.Zero) {
		movementType = MovementTypeAdjustIn
	} else if difference.LessThan(decimal.Zero) {
		movementType = MovementTypeAdjustOut
		difference = difference.Abs()
	} else {
		// Não há diferença, não criar movimento
		return nil, nil
	}

	movement := &StockMovement{
		Entity: entity.NewEntity(),
		StockMovementCommonAttributes: StockMovementCommonAttributes{
			StockID:    s.ID,
			Type:       movementType,
			Quantity:   difference,
			Reason:     reason,
			EmployeeID: employeeID,
		},
	}

	return movement, nil
}

// ReserveStock reserva estoque para um pedido (quando pedido fica pending)
func (s *Stock) ReserveStock(quantity decimal.Decimal, orderID uuid.UUID, employeeID uuid.UUID, price decimal.Decimal, totalPrice decimal.Decimal) (*StockMovement, error) {
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
			Type:       MovementTypeOut,
			Quantity:   quantity,
			Reason:     "Reserva automática - Pedido pendente",
			EmployeeID: employeeID,
			OrderID:    &orderID,
			Price:      price,
			TotalPrice: totalPrice,
		},
	}

	return movement, nil
}

// RestoreStock restaura estoque quando pedido é cancelado
func (s *Stock) RestoreStock(quantity decimal.Decimal, orderID uuid.UUID, employeeID uuid.UUID, price decimal.Decimal, totalPrice decimal.Decimal) (*StockMovement, error) {
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
			Type:       MovementTypeIn,
			Quantity:   quantity,
			Reason:     "Restauração automática - Pedido cancelado",
			EmployeeID: employeeID,
			OrderID:    &orderID,
			Price:      price,
			TotalPrice: totalPrice,
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
				StockID:     s.ID,
				Type:        AlertTypeLowStock,
				Message:     "Estoque abaixo do mínimo",
				IsResolved:  false,
				ProductName: s.Product.Name,
				ProductSKU:  s.Product.SKU,
			},
		})
	}

	// Alerta de estoque zerado
	if s.CurrentStock.LessThanOrEqual(decimal.Zero) {
		alerts = append(alerts, &StockAlert{
			Entity: entity.NewEntity(),
			StockAlertCommonAttributes: StockAlertCommonAttributes{
				StockID:     s.ID,
				Type:        AlertTypeOutOfStock,
				Message:     "Produto sem estoque",
				IsResolved:  false,
				ProductName: s.Product.Name,
				ProductSKU:  s.Product.SKU,
			},
		})
	}

	// Alerta de estoque acima do máximo
	if s.MaxStock.GreaterThan(decimal.Zero) && s.CurrentStock.GreaterThan(s.MaxStock) {
		alerts = append(alerts, &StockAlert{
			Entity: entity.NewEntity(),
			StockAlertCommonAttributes: StockAlertCommonAttributes{
				StockID:     s.ID,
				Type:        AlertTypeOverStock,
				Message:     "Estoque acima do máximo",
				IsResolved:  false,
				ProductName: s.Product.Name,
				ProductSKU:  s.Product.SKU,
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
