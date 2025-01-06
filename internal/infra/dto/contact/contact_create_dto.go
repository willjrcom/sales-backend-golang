package contactdto

import (
	"errors"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrDddIsEmpty    = errors.New("ddd is required")
	ErrNumberIsEmpty = errors.New("number is required")
)

type ContactCreateDTO struct {
	Ddd    string                   `json:"ddd"`
	Number string                   `json:"number"`
	Type   personentity.ContactType `json:"type"`
}

func (c *ContactCreateDTO) validate() error {
	if c.Ddd == "" {
		return ErrDddIsEmpty
	}

	if c.Number == "" {
		return ErrNumberIsEmpty
	}

	return nil
}

func (c *ContactCreateDTO) ToModel() (*personentity.Contact, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	contactCommonAttributes := &personentity.ContactCommonAttributes{
		Ddd:    c.Ddd,
		Number: c.Number,
		Type:   c.Type,
	}

	return personentity.NewContact(contactCommonAttributes), nil
}
