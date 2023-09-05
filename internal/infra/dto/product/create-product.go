package productdto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CreateProductInput struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	Category    string  `json:"category"`
	IsAvailable bool    `json:"is_available"`
}

func (p *CreateProductInput) Validate() error {
	if p.Price > p.Cost {
		return errors.New("price must be greater than cost")
	}

	if p.Category == "" {
		return errors.New("category is required")
	}

	return nil
}

func (p *CreateProductInput) ToModel() (*productentity.Product, error) {
	if err := p.Validate(); err != nil {
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
		Category:    productentity.CategoryProduct{Name: p.Category},
		IsAvailable: p.IsAvailable,
	}, nil
}
