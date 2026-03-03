package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProductVariation struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:product_variations,alias:product_variation"`
	ProductID     uuid.UUID        `bun:"product_id,type:uuid,notnull"`
	SizeID        uuid.UUID        `bun:"size_id,type:uuid,notnull"`
	Size          *Size            `bun:"rel:belongs-to"`
	Price         *decimal.Decimal `bun:"price,type:decimal(10,2),notnull"`
	IsAvailable   bool             `bun:"is_available,notnull,default:true"`
}

// ToDomain converte model para domain
func (pv *ProductVariation) ToDomain() productentity.ProductVariation {
	variation := productentity.ProductVariation{
		Entity:      pv.Entity.ToDomain(),
		ProductID:   pv.ProductID,
		SizeID:      pv.SizeID,
		Price:       pv.GetPrice(),
		IsAvailable: pv.IsAvailable,
	}

	if pv.Size != nil {
		variation.Size = pv.Size.ToDomain()
	}

	return variation
}

// FromDomain converte domain para model
func (pv *ProductVariation) FromDomain(variation productentity.ProductVariation) {
	pv.ID = variation.ID
	pv.ProductID = variation.ProductID
	pv.SizeID = variation.SizeID
	pv.Price = &variation.Price
	pv.IsAvailable = variation.IsAvailable

	if variation.Size != nil {
		pv.Size = &Size{}
		pv.Size.FromDomain(variation.Size)
	}
}

func (p *ProductVariation) GetPrice() decimal.Decimal {
	if p.Price == nil {
		return decimal.Zero
	}
	return *p.Price
}
