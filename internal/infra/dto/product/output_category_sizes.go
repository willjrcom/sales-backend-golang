package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategorySizesOutput struct {
	ID    uuid.UUID    `json:"id"`
	Name  string       `json:"name"`
	Sizes []SizeOutput `json:"sizes"`
}

func (c *CategorySizesOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name

	c.Sizes = make([]SizeOutput, len(model.Sizes))

	for i, v := range model.Sizes {
		c.Sizes[i].FromModel(&v)
	}

}
