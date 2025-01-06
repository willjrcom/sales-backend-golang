package contactdto

import (
	"github.com/google/uuid"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type ContactDTO struct {
	ID     uuid.UUID                `json:"id"`
	Ddd    string                   `json:"ddd"`
	Number string                   `json:"number"`
	Type   personentity.ContactType `json:"type"`
}

func (c *ContactDTO) FromDomain(contact *personentity.Contact) {
	*c = ContactDTO{
		ID:     contact.ID,
		Ddd:    contact.Ddd,
		Number: contact.Number,
		Type:   contact.Type,
	}
}
