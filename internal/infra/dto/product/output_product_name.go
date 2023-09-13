package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductNameOutput struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
}

func (p *ProductNameOutput) FromModel(model *productentity.Product) {
	p.ID = model.ID
	p.Code = model.Code
	p.Name = model.Name
	p.Price = model.Price
	p.IsAvailable = model.IsAvailable
}
