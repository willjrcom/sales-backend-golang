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
	Price       decimal.Decimal  `json:"price"`
	Cost        decimal.Decimal  `json:"cost"`
	IsAvailable bool             `json:"is_available"`
}

func (d *ProductVariationDTO) FromDomain(variation productentity.ProductVariation) {
	*d = ProductVariationDTO{
		ID:          variation.ID,
		ProductID:   variation.ProductID,
		SizeID:      variation.SizeID,
		Price:       variation.Price,
		Cost:        variation.Cost,
		IsAvailable: variation.IsAvailable,
	}

	if variation.Size != nil {
		d.Size = &sizedto.SizeDTO{}
		d.Size.FromDomain(variation.Size)
	}
}

type ProductVariationCreateDTO struct {
	SizeID      uuid.UUID       `json:"size_id"`
	Price       decimal.Decimal `json:"price"`
	Cost        decimal.Decimal `json:"cost"`
	IsAvailable bool            `json:"is_available"`
}

func (d *ProductVariationCreateDTO) ToDomain() productentity.ProductVariation {
	return productentity.ProductVariation{
		Entity:      entity.NewEntity(),
		SizeID:      d.SizeID,
		Price:       d.Price,
		Cost:        d.Cost,
		IsAvailable: d.IsAvailable,
	}
}
