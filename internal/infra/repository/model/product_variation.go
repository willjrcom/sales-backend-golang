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
	Cost          *decimal.Decimal `bun:"cost,type:decimal(10,2)"`
	IsAvailable   bool             `bun:"is_available"`
}

func (p *ProductVariation) ToDomain() productentity.ProductVariation {
	return productentity.ProductVariation{
		Entity:      p.Entity.ToDomain(),
		ProductID:   p.ProductID,
		SizeID:      p.SizeID,
		Size:        p.Size.ToDomain(),
		Price:       p.GetPrice(),
		Cost:        p.GetCost(),
		IsAvailable: p.IsAvailable,
	}
}

func (p *ProductVariation) FromDomain(productVariation productentity.ProductVariation) {
	*p = ProductVariation{
		Entity:      entitymodel.FromDomain(productVariation.Entity),
		ProductID:   productVariation.ProductID,
		SizeID:      productVariation.SizeID,
		Price:       &productVariation.Price,
		Cost:        &productVariation.Cost,
		IsAvailable: productVariation.IsAvailable,
	}

	if productVariation.Size != nil {
		p.Size = &Size{}
		p.Size.FromDomain(productVariation.Size)
	}
}

func (p *ProductVariation) GetPrice() decimal.Decimal {
	if p.Price == nil {
		return decimal.Zero
	}
	return *p.Price
}

func (p *ProductVariation) GetCost() decimal.Decimal {
	if p.Cost == nil {
		return decimal.Zero
	}
	return *p.Cost
}
