package productdto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrCodeRequired         = errors.New("code is required")
	ErrNameRequired         = errors.New("name is required")
	ErrCostGreaterThanPrice = errors.New("Cost must be greater than Price")
	ErrCategoryRequired     = errors.New("category is required")
)

type CreateProductInput struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Size        string    `json:"size"`
	Price       float64   `json:"price"`
	Cost        float64   `json:"cost"`
	CategoryID  uuid.UUID `json:"category_id"`
	IsAvailable bool      `json:"is_available"`
}

func (p *CreateProductInput) validate() error {
	if p.Code == "" {
		return ErrCodeRequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price < p.Cost {
		return ErrCostGreaterThanPrice
	}
	if p.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (p *CreateProductInput) ToModel() (*productentity.Product, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	return &productentity.Product{
		Entity:      entity.Entity{ID: uuid.New(), CreatedAt: time.Now()},
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		Size:        p.Size,
		Price:       p.Price,
		Cost:        p.Cost,
		CategoryID:  p.CategoryID,
		IsAvailable: p.IsAvailable,
	}, nil
}
