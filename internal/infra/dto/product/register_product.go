package productdto

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrCodeRequired         = errors.New("code is required")
	ErrNameRequired         = errors.New("name is required")
	ErrCostGreaterThanPrice = errors.New("cost must be greater than Price")
	ErrCategoryRequired     = errors.New("category is required")
	ErrSizeRequired         = errors.New("size is required")
)

type RegisterProductInput struct {
	productentity.ProductCommonAttributes
	Image *multipart.File `json:"image"`
}

func (p *RegisterProductInput) validate() error {
	if p.Code == "" {
		return ErrCodeRequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price < p.Cost {
		return ErrCostGreaterThanPrice
	}
	if len(p.CategoryID.String()) == 0 || p.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}
	if p.SizeID == uuid.Nil {
		return ErrSizeRequired
	}

	return nil
}

func (p *RegisterProductInput) ToModel() (*productentity.Product, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	productCommonAttributes := productentity.ProductCommonAttributes{
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		SizeID:      p.SizeID,
		Price:       p.Price,
		Cost:        p.Cost,
		CategoryID:  p.CategoryID,
		IsAvailable: p.IsAvailable,
	}

	return &productentity.Product{
		Entity:                  entity.NewEntity(),
		ProductCommonAttributes: productCommonAttributes,
	}, nil
}
