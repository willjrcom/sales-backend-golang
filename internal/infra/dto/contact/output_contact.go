package contactdto

import (
	"github.com/google/uuid"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type ContactOutput struct {
	ID uuid.UUID `json:"id"`
	personentity.ContactCommonAttributes
}

func (c *ContactOutput) FromModel(model *personentity.Contact) {
	c.ID = model.ID
	c.ContactCommonAttributes = model.ContactCommonAttributes
}
