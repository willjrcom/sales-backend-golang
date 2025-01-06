package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductCategoryToComplement struct {
	bun.BaseModel        `bun:"table:product_category_to_complement"`
	CategoryID           uuid.UUID        `bun:"type:uuid,pk"`
	Category             *ProductCategory `bun:"rel:belongs-to,join:category_id=id"`
	ComplementCategoryID uuid.UUID        `bun:"type:uuid,pk"`
	ComplementCategory   *ProductCategory `bun:"rel:belongs-to,join:complement_category_id=id"`
}

func (p *ProductCategoryToComplement) FromDomain(productCategoryToComplement *productentity.ProductCategoryToComplement) {
	*p = ProductCategoryToComplement{
		CategoryID:           productCategoryToComplement.CategoryID,
		ComplementCategoryID: productCategoryToComplement.ComplementCategoryID,
	}
}

func (p *ProductCategoryToComplement) ToDomain() *productentity.ProductCategoryToComplement {
	if p == nil {
		return nil
	}
	return &productentity.ProductCategoryToComplement{
		CategoryID:           p.CategoryID,
		ComplementCategoryID: p.ComplementCategoryID,
	}
}
