package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryProductsOutput struct {
	ID       uuid.UUID           `json:"id"`
	Name     string              `json:"name"`
	Products []ProductNameOutput `json:"products"`
}

func (c *CategoryProductsOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name
	c.Products = make([]ProductNameOutput, len(model.Products))

	for i, v := range model.Products {
		c.Products[i].FromModel(&v)
	}
}
