package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Size struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:sizes"`
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string    `bun:"name"`
	IsActive   *bool     `bun:"is_active"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	Products   []Product `bun:"rel:has-many,join:id=size_id"`
}

func (s *Size) FromDomain(size *productentity.Size) {
	if size == nil {
		return
	}
	*s = Size{
		Entity: entitymodel.FromDomain(size.Entity),
		SizeCommonAttributes: SizeCommonAttributes{
			Name:       size.Name,
			IsActive:   size.IsActive,
			CategoryID: size.CategoryID,
		},
	}

	for _, product := range size.Products {
		p := Product{}
		p.FromDomain(&product)
		s.Products = append(s.Products, p)
	}
}

func (s *Size) ToDomain() *productentity.Size {
	if s == nil {
		return nil
	}
	size := &productentity.Size{
		Entity: s.Entity.ToDomain(),
		SizeCommonAttributes: productentity.SizeCommonAttributes{
			Name:       s.Name,
			IsActive:   s.IsActive,
			CategoryID: s.CategoryID,
		},
	}

	for _, product := range s.Products {
		size.Products = append(size.Products, *product.ToDomain())
	}

	return size
}
