package contactdto

import (
	"github.com/google/uuid"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type ContactOutput struct {
	ID       uuid.UUID `json:"id"`
	PersonID uuid.UUID `json:"person_id,omitempty"`
	Ddd      string    `json:"ddd"`
	Number   string    `json:"number"`
}

func (c *ContactOutput) FromModel(model *personentity.Contact) {
	c.ID = model.ID
	c.Ddd = model.Ddd
	c.Number = model.Number
	c.PersonID = model.PersonID
}
