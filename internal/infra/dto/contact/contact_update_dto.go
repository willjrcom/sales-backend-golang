package contactdto

import personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"

type ContactUpdateDTO struct {
	Number string `json:"number"`
}

func (c *ContactUpdateDTO) validate() error {
	return nil
}

func (c *ContactUpdateDTO) UpdateDomain(contact *personentity.Contact, Type personentity.ContactType) error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.Number != "" {
		contact.Number = c.Number
	}

	contact.Type = Type

	return nil
}
