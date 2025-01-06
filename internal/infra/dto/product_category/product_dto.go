package productcategorydto

import (
	"github.com/google/uuid"
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
	Price       float64          `json:"price"`
	Cost        float64          `json:"cost"`
	IsAvailable bool             `json:"is_available"`
	CategoryID  uuid.UUID        `json:"category_id"`
	Category    *CategoryDTO     `json:"category"`
	SizeID      uuid.UUID        `json:"size_id"`
	Size        *sizedto.SizeDTO `json:"size"`
}

func FromDomain(product *productentity.Product) *ProductDTO {
	p := &ProductDTO{
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
		SizeID:      product.SizeID,
	}

	p.Category.FromDomain(product.Category)
	p.Size.FromDomain(product.Size)

	return p
}
