package contactdto

import (
	"errors"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidContent = errors.New("invalid content")
)

type UpdateContactInput struct {
	Contact string `json:"contact"`
}

func (c *UpdateContactInput) validate() error {
	if c.Contact == "" {
		return ErrInvalidContent
	}

	return nil
}

func (c *UpdateContactInput) UpdateModel(model *personentity.Contact) error {
	if err := c.validate(); err != nil {
		return err
	}

	ddd, number, err := personentity.ValidateAndExtractContact(c.Contact)

	if err != nil {
		return err
	}

	model.Ddd = ddd
	model.Number = number
	return nil
}
