package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Size struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:sizes,alias:size"`
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string              `bun:"name"`
	IsActive   bool                `bun:"is_active"`
	CategoryID uuid.UUID           `bun:"column:category_id,type:uuid,notnull"`
	Variations []*ProductVariation `bun:"rel:has-many,join:id=size_id"`
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

	for _, variation := range size.Variations {
		v := &ProductVariation{}
		v.FromDomain(variation)
		s.Variations = append(s.Variations, v)
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

	for _, variation := range s.Variations {
		size.Variations = append(size.Variations, variation.ToDomain())
	}

	return size
}
