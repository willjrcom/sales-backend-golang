package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	categorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/category"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size_category"
)

type ProductOutput struct {
	ID          uuid.UUID                       `json:"id"`
	Code        string                          `json:"code"`
	Name        string                          `json:"name"`
	Description string                          `json:"description"`
	Price       float64                         `json:"price"`
	IsAvailable bool                            `json:"is_available"`
	Category    *categorydto.CategoryNameOutput `json:"category,omitempty"`
	Size        *sizedto.SizeNameOutput         `json:"size,omitempty"`
}

func (p *ProductOutput) FromModel(model *productentity.Product) {
	p.ID = model.ID
	p.Code = model.Code
	p.Name = model.Name
	p.Description = model.Description
	p.Price = model.Price
	p.IsAvailable = model.IsAvailable

	if model.Category != nil {
		p.Category = &categorydto.CategoryNameOutput{}
		p.Category.FromModel(model.Category)
	}

	if model.Size != nil {
		p.Size = &sizedto.SizeNameOutput{}
		p.Size.FromModel(model.Size)
	}
}
