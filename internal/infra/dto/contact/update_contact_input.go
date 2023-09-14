package contactdto

import (
	"errors"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidContent = errors.New("invalid content")
)

type UpdateContactInput struct {
	Ddd    *string `json:"ddd"`
	Number *string `json:"number"`
}

func (c *UpdateContactInput) validate() error {
	if c.Ddd == nil && c.Number == nil {
		return ErrInvalidContent
	}

	return nil
}

func (c *UpdateContactInput) UpdateModel(model *personentity.Contact) error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.Ddd != nil {
		model.Ddd = *c.Ddd
	}
	if c.Number != nil {
		model.Number = *c.Number
	}

	return nil
}
