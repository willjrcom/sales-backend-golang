package productcategorydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
)

type ProductDTO struct {
	ID          uuid.UUID        `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Flavors     []string         `json:"flavors"`
	ImagePath   *string          `json:"image_path"`
	Description string           `json:"description"`
	Price       decimal.Decimal  `json:"price"`
	Cost        decimal.Decimal  `json:"cost"`
	IsAvailable bool             `json:"is_available"`
	CategoryID  uuid.UUID        `json:"category_id"`
	Category    *CategoryDTO     `json:"category"`
	SizeID      uuid.UUID        `json:"size_id"`
	Size        *sizedto.SizeDTO `json:"size"`
}

func (p *ProductDTO) FromDomain(product *productentity.Product) {
	if product == nil {
		return
	}

	*p = ProductDTO{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		Flavors:     product.Flavors,
		ImagePath:   product.ImagePath,
		Description: product.Description,
		Price:       product.Price,
		Cost:        product.Cost,
		IsAvailable: product.IsAvailable,
		CategoryID:  product.CategoryID,
		Category:    &CategoryDTO{},
		SizeID:      product.SizeID,
		Size:        &sizedto.SizeDTO{},
	}

	p.Category.FromDomain(product.Category)
	p.Size.FromDomain(product.Size)

	if product.Category == nil {
		p.Category = nil
	}
	if product.Size == nil {
		p.Size = nil
	}
}
