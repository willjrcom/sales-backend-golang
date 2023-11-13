package categorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity_category"
)

type CategoryQuantitiesOutput struct {
	ID         uuid.UUID                    `json:"id"`
	Name       string                       `json:"name"`
	Quantities []quantitydto.QuantityOutput `json:"sizes"`
}

func (c *CategoryQuantitiesOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name

	c.Quantities = make([]quantitydto.QuantityOutput, len(model.Quantities))

	for i, v := range model.Quantities {
		c.Quantities[i].FromModel(&v)
	}

}
