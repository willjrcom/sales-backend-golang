package contactdto

import (
	"github.com/google/uuid"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type ContactoDTO struct {
	ID     uuid.UUID                `json:"id"`
	Ddd    string                   `json:"ddd"`
	Number string                   `json:"number"`
	Type   personentity.ContactType `json:"type"`
}

func (c *ContactoDTO) FromDomain(contact *personentity.Contact) {
	*c = ContactoDTO{
		ID:     contact.ID,
		Ddd:    contact.Ddd,
		Number: contact.Number,
		Type:   contact.Type,
	}
}
