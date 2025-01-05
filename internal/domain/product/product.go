package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("product category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	bun.BaseModel `bun:"table:products"`
	ProductCommonAttributes
}

type ProductCommonAttributes struct {
	Code        string           `bun:"code,notnull" json:"code"`
	Name        string           `bun:"name,notnull" json:"name"`
	Flavors     []string         `bun:"flavors,type:jsonb" json:"flavors,omitempty"`
	ImagePath   *string          `bun:"image_path" json:"image_path"`
	Description string           `bun:"description" json:"description"`
	Price       float64          `bun:"price,notnull" json:"price"`
	Cost        float64          `bun:"cost" json:"cost"`
	IsAvailable bool             `bun:"is_available" json:"is_available"`
	CategoryID  uuid.UUID        `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category    *ProductCategory `bun:"rel:belongs-to" json:"category,omitempty"`
	SizeID      uuid.UUID        `bun:"column:size_id,type:uuid,notnull" json:"size_id"`
	Size        *Size            `bun:"rel:belongs-to" json:"size,omitempty"`
}

type PatchProduct struct {
	Code        *string    `json:"code"`
	Name        *string    `json:"name"`
	Flavors     []string   `json:"flavors,omitempty"`
	ImagePath   *string    `json:"image_path"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price"`
	Cost        *float64   `json:"cost"`
	IsAvailable *bool      `json:"is_available"`
	CategoryID  *uuid.UUID `json:"category_id"`
	SizeID      *uuid.UUID `json:"size_id"`
}

func NewProduct(productCommonAttributes ProductCommonAttributes) *Product {
	return &Product{
		Entity:                  entity.NewEntity(),
		ProductCommonAttributes: productCommonAttributes,
	}
}
func UpdateProduct(patchProduct PatchProduct) *Product {
	productCommonAttributes := ProductCommonAttributes{}

	if patchProduct.Code != nil {
		productCommonAttributes.Code = *patchProduct.Code
	}
	if patchProduct.Name != nil {
		productCommonAttributes.Name = *patchProduct.Name
	}
	if patchProduct.Flavors != nil {
		productCommonAttributes.Flavors = patchProduct.Flavors
	}
	if patchProduct.ImagePath != nil {
		productCommonAttributes.ImagePath = patchProduct.ImagePath
	}
	if patchProduct.Description != nil {
		productCommonAttributes.Description = *patchProduct.Description
	}
	if patchProduct.Price != nil {
		productCommonAttributes.Price = *patchProduct.Price
	}
	if patchProduct.Cost != nil {
		productCommonAttributes.Cost = *patchProduct.Cost
	}
	if patchProduct.IsAvailable != nil {
		productCommonAttributes.IsAvailable = *patchProduct.IsAvailable
	}
	if patchProduct.CategoryID != nil {
		productCommonAttributes.CategoryID = *patchProduct.CategoryID
	}
	if patchProduct.SizeID != nil {
		productCommonAttributes.SizeID = *patchProduct.SizeID
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
