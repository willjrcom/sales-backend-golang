package productcategorydto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	Price       decimal.Decimal `json:"price"`
	Cost        decimal.Decimal `json:"cost"`
	IsAvailable bool            `json:"is_available"`
	IsActive    *bool           `json:"is_active"`
	CategoryID  *uuid.UUID      `json:"category_id"`
	SizeID      *uuid.UUID      `json:"size_id"`
	ImagePath   string          `json:"image_path"`
}

func (p *ProductCreateDTO) validate() error {
	if p.Code == "" {
		return ErrCodeRequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price.LessThan(p.Cost) {
		return ErrCostGreaterThanPrice
	}
	if p.CategoryID == nil || len(p.CategoryID.String()) == 0 {
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

	flavors := p.Flavors
	if len(flavors) == 0 {
		flavors = []string{}
	}

	isActive := true
	if p.IsActive != nil {
		isActive = *p.IsActive
	}

	productCommonAttributes := productentity.ProductCommonAttributes{
		Code:        p.Code,
		Name:        p.Name,
		Flavors:     flavors,
		Description: p.Description,
		Price:       p.Price,
		Cost:        p.Cost,
		IsAvailable: p.IsAvailable,
		IsActive:    isActive,
		CategoryID:  *p.CategoryID,
		SizeID:      *p.SizeID,
		ImagePath:   &p.ImagePath,
	}

	return productentity.NewProduct(productCommonAttributes), nil
}
