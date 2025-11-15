package stockdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

// StockWithProductDTO representa o DTO de estoque com informações do produto
type StockWithProductDTO struct {
	Stock   *StockDTO                      `json:"stock,omitempty"`
	Product *productcategorydto.ProductDTO `json:"product,omitempty"`
}

func (s *StockWithProductDTO) FromDomain(stock *stockentity.Stock, product *productentity.Product) {
	if stock != nil {
		stockDTO := &StockDTO{}
		stockDTO.FromDomain(stock)
		s.Stock = stockDTO
	}

	if product != nil {
		productDTO := productcategorydto.ProductDTO{}
		productDTO.FromDomain(product)
		s.Product = &productDTO
	}
}
