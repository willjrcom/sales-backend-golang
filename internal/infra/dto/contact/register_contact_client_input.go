package contactdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrContactIsEmpty  = errors.New("contact is required")
	ErrPersonIdIsEmpty = errors.New("person ID is required")
)

type RegisterContactClientInput struct {
	PersonID uuid.UUID `json:"person_id"`
	Contact  string    `json:"contact"`
}

func (c *RegisterContactClientInput) validate() error {
	if c.Contact == "" {
		return ErrContactIsEmpty
	}
	if c.PersonID == uuid.Nil {
		return ErrPersonIdIsEmpty
	}

	return nil
}

func (c *RegisterContactClientInput) ToModel() (*personentity.Contact, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	// Validate contact
	ddd, number, err := personentity.ValidateAndExtractContact(c.Contact)

	if err != nil {
		return nil, err
	}

	return &personentity.Contact{
		Entity:   entity.NewEntity(),
		Ddd:      ddd,
		Number:   number,
		PersonID: c.PersonID,
		Type:     "client",
	}, nil
}
