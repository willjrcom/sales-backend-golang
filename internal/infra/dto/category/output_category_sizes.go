package categorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size_category"
)

type CategorySizesOutput struct {
	ID    uuid.UUID            `json:"id"`
	Name  string               `json:"name"`
	Sizes []sizedto.SizeOutput `json:"sizes"`
}

func (c *CategorySizesOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name

	c.Sizes = make([]sizedto.SizeOutput, len(model.Sizes))

	for i, v := range model.Sizes {
		c.Sizes[i].FromModel(&v)
	}

}
