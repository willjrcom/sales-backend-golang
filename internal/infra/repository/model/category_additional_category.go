package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductCategoryToAdditional struct {
	bun.BaseModel        `bun:"table:product_category_to_additional"`
	CategoryID           uuid.UUID        `bun:"type:uuid,pk"`
	Category             *ProductCategory `bun:"rel:belongs-to,join:category_id=id"`
	AdditionalCategoryID uuid.UUID        `bun:"type:uuid,pk"`
	AdditionalCategory   *ProductCategory `bun:"rel:belongs-to,join:additional_category_id=id"`
}

func (p *ProductCategoryToAdditional) FromDomain(productCategoryToAdditional *productentity.ProductCategoryToAdditional) {
	*p = ProductCategoryToAdditional{
		CategoryID:           productCategoryToAdditional.CategoryID,
		AdditionalCategoryID: productCategoryToAdditional.AdditionalCategoryID,
	}
}

func (p *ProductCategoryToAdditional) ToDomain() *productentity.ProductCategoryToAdditional {
	return &productentity.ProductCategoryToAdditional{
		CategoryID:           p.CategoryID,
		AdditionalCategoryID: p.AdditionalCategoryID,
	}
}
