package categorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type CategoryNameOutput struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (c *CategoryNameOutput) FromModel(model *productentity.Category) {
	c.ID = model.ID
	c.Name = model.Name
}
