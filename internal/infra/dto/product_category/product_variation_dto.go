package productcategorydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
)

type ProductVariationDTO struct {
	ID          uuid.UUID        `json:"id"`
	ProductID   uuid.UUID        `json:"product_id"`
	SizeID      uuid.UUID        `json:"size_id"`
	Size        *sizedto.SizeDTO `json:"size"`
	Price       *decimal.Decimal `json:"price"`
	IsAvailable bool             `json:"is_available"`
}

func (pv *ProductVariationDTO) ToDomain() productentity.ProductVariation {
	variation := productentity.ProductVariation{
		Entity:      entity.Entity{ID: pv.ID},
		ProductID:   pv.ProductID,
		SizeID:      pv.SizeID,
		Price:       *pv.Price, // Assuming Price is always non-nil when converting to domain
		IsAvailable: pv.IsAvailable,
	}

	return variation
}

func (pv *ProductVariationDTO) FromDomain(variation productentity.ProductVariation) {
	pv.ID = variation.ID
	pv.ProductID = variation.ProductID
	pv.SizeID = variation.SizeID
	pv.Price = &variation.Price
	pv.IsAvailable = variation.IsAvailable

	if variation.Size != nil {
		pv.Size = &sizedto.SizeDTO{}
		pv.Size.FromDomain(variation.Size)
	}
}

type ProductVariationCreateDTO struct {
	SizeID      uuid.UUID       `json:"size_id"`
	Price       decimal.Decimal `json:"price"`
	IsAvailable bool            `json:"is_available"`
}

func (d *ProductVariationCreateDTO) ToDomain(productID uuid.UUID) productentity.ProductVariation {
	return *productentity.NewProductVariation(productID, d.SizeID, d.Price)
}
