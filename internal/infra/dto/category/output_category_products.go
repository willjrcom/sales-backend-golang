package categorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryProductsOutput struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Products []string  `json:"product_ids"`
}

func (c *CategoryProductsOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name
	c.Products = make([]string, len(model.Products))

	for i, v := range model.Products {
		c.Products[i] = v.ID.String()
	}
}
