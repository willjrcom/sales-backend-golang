package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("product category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	ProductCommonAttributes
	Variations []ProductVariation
}

type ProductCommonAttributes struct {
	SKU         string
	Name        string
	Flavors     []string
	ImagePath   *string
	Description string
	IsActive    bool
	CategoryID  uuid.UUID
	Category    *ProductCategory
}

func (p *Product) AddVariation(variation ProductVariation) {
	p.Variations = append(p.Variations, variation)
}

func NewProduct(productCommonAttributes ProductCommonAttributes) *Product {
	if len(productCommonAttributes.Flavors) == 0 {
		productCommonAttributes.Flavors = []string{}
	}

	return &Product{
		Entity:                  entity.NewEntity(),
		ProductCommonAttributes: productCommonAttributes,
		Variations:              []ProductVariation{},
	}
}
