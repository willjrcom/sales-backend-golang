package stockdto

import (
	"time"

	"github.com/google/uuid"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

// StockAlertDTO representa o DTO de alerta de estoque
type StockAlertDTO struct {
	ID                 uuid.UUID                               `json:"id"`
	StockID            uuid.UUID                               `json:"stock_id"`
	Type               string                                  `json:"type"`
	Message            string                                  `json:"message"`
	IsResolved         bool                                    `json:"is_resolved"`
	ResolvedAt         *time.Time                              `json:"resolved_at,omitempty"`
	ResolvedBy         *uuid.UUID                              `json:"resolved_by,omitempty"`
	CreatedAt          time.Time                               `json:"created_at"`
	ProductID          uuid.UUID                               `json:"product_id"`
	Product            *productcategorydto.ProductDTO          `json:"product"`
	ProductVariationID uuid.UUID                               `json:"product_variation_id"`
	ProductVariation   *productcategorydto.ProductVariationDTO `json:"product_variation"`
}

// FromDomain converte domain para DTO
func (sa *StockAlertDTO) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	var productDTO *productcategorydto.ProductDTO
	if alert.Product != nil {
		productDTO = &productcategorydto.ProductDTO{}
		productDTO.FromDomain(alert.Product)
	}

	var variationDTO *productcategorydto.ProductVariationDTO
	if alert.ProductVariation != nil {
		variationDTO = &productcategorydto.ProductVariationDTO{}
		variationDTO.FromDomain(*alert.ProductVariation)
	}

	*sa = StockAlertDTO{
		ID:                 alert.ID,
		StockID:            alert.StockID,
		Type:               string(alert.Type),
		Message:            alert.Message,
		IsResolved:         alert.IsResolved,
		ResolvedAt:         alert.ResolvedAt,
		ResolvedBy:         alert.ResolvedBy,
		CreatedAt:          alert.CreatedAt,
		ProductID:          alert.ProductID,
		Product:            productDTO,
		ProductVariationID: alert.ProductVariationID,
		ProductVariation:   variationDTO,
	}
}
