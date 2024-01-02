package clientdto

import (
	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type ClientOutput struct {
	ID uuid.UUID `json:"id"`
	personentity.PersonCommonAttributes
}

func (c *ClientOutput) FromModel(model *cliententity.Client) {
	c.ID = model.ID
	c.PersonCommonAttributes = model.PersonCommonAttributes
}
