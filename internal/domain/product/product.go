package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("product category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	ProductCommonAttributes
}

type ProductCommonAttributes struct {
	Code        string
	Name        string
	Flavors     []string
	ImagePath   *string
	Description string
	Price       decimal.Decimal
	Cost        decimal.Decimal
	IsAvailable bool
	IsActive    bool
	CategoryID  uuid.UUID
	Category    *ProductCategory
	SizeID      uuid.UUID
	Size        *Size
}

func NewProduct(productCommonAttributes ProductCommonAttributes) *Product {
	if len(productCommonAttributes.Flavors) == 0 {
		productCommonAttributes.Flavors = []string{}
	}

	return &Product{
		Entity:                  entity.NewEntity(),
		ProductCommonAttributes: productCommonAttributes,
	}
}

func (p *Product) FindSizeInCategory() (bool, error) {
	if p.Category == nil {
		return false, ErrCategoryNotFound
	}

	for _, v := range p.Category.Sizes {
		if v.ID == p.SizeID {
			return true, nil
		}
	}

	return false, errors.New("size not found")
}
