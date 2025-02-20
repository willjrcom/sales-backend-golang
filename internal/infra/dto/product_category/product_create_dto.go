package productcategorydto

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrCodeRequired         = errors.New("code is required")
	ErrNameRequired         = errors.New("name is required")
	ErrCostGreaterThanPrice = errors.New("cost must be greater than Price")
	ErrCategoryRequired     = errors.New("category is required")
	ErrSizeRequired         = errors.New("size is required")
)

type ProductCreateDTO struct {
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Flavors     []string        `json:"flavors"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Cost        float64         `json:"cost"`
	IsAvailable bool            `json:"is_available"`
	CategoryID  *uuid.UUID      `json:"category_id"`
	SizeID      *uuid.UUID      `json:"size_id"`
	Image       *multipart.File `json:"image"`
}

func (p *ProductCreateDTO) validate() error {
	if p.Code == "" {
		return ErrCodeRequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price > 0 && p.Cost > 0 && p.Price < p.Cost {
		return ErrCostGreaterThanPrice
	}
	if len(p.CategoryID.String()) == 0 || p.CategoryID == nil {
		return ErrCategoryRequired
	}
	if p.SizeID == nil {
		return ErrSizeRequired
	}

	return nil
}

func (p *ProductCreateDTO) ToDomain() (*productentity.Product, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	productCommonAttributes := productentity.ProductCommonAttributes{
		Code:        p.Code,
		Name:        p.Name,
		Flavors:     p.Flavors,
		Description: p.Description,
		Price:       p.Price,
		Cost:        p.Cost,
		IsAvailable: p.IsAvailable,
		CategoryID:  *p.CategoryID,
		SizeID:      *p.SizeID,
	}

	return productentity.NewProduct(productCommonAttributes), nil
}
