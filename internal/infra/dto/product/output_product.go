package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProductOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.ProductCommonAttributes
}

func (p *ProductOutput) FromModel(model *productentity.Product) {
	p.ID = model.ID
	p.ProductCommonAttributes = model.ProductCommonAttributes
}
