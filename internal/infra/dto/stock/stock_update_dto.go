package stockdto

import (
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
)

type StockUpdateDTO struct {
	MinStock *decimal.Decimal `json:"min_stock,omitempty"`
	MaxStock *decimal.Decimal `json:"max_stock,omitempty"`
	Unit     *string          `json:"unit,omitempty"`
	IsActive *bool            `json:"is_active,omitempty"`
}

func (s *StockUpdateDTO) Validate() error {
	return nil
}

func (s *StockUpdateDTO) UpdateDomain(stock *stockentity.Stock) (err error) {
	if err := s.Validate(); err != nil {
		return err
	}

	if s.MinStock != nil {
		stock.MinStock = *s.MinStock
	}
	if s.MaxStock != nil {
		stock.MaxStock = *s.MaxStock
	}
	if s.Unit != nil {
		stock.Unit = *s.Unit
	}
	if s.IsActive != nil {
		stock.IsActive = *s.IsActive
	}

	return nil
}
