package contactdto

import (
	"errors"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrNumberIsEmpty = errors.New("number is required")
)

type ContactCreateDTO struct {
	Number string                   `json:"number"`
	Type   personentity.ContactType `json:"type"`
}

func (c *ContactCreateDTO) validate() error {
	if c.Number == "" {
		return ErrNumberIsEmpty
	}

	return nil
}

func (c *ContactCreateDTO) ToDomain() (*personentity.Contact, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	contactCommonAttributes := &personentity.ContactCommonAttributes{
		Number: c.Number,
		Type:   c.Type,
	}

	return personentity.NewContact(contactCommonAttributes), nil
}
