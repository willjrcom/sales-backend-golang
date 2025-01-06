package contactdto

import personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"

type ContactUpdateDTO struct {
	Ddd    *string `json:"ddd"`
	Number *string `json:"number"`
}

func (c *ContactUpdateDTO) validate() error {
	return nil
}

func (c *ContactUpdateDTO) UpdateDomain(contact *personentity.Contact) error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.Ddd != nil {
		contact.Ddd = *c.Ddd
	}
	if c.Number != nil {
		contact.Number = *c.Number
	}

	return nil
}
