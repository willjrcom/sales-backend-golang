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
	ClientID uuid.UUID `json:"client_id"`
	Contact  string    `json:"contact"`
}

func (c *RegisterContactClientInput) validate() error {
	if c.Contact == "" {
		return ErrContactIsEmpty
	}
	if c.ClientID == uuid.Nil {
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

	contactCommonAttributes := personentity.ContactCommonAttributes{
		PersonID: c.ClientID,
		Ddd:      ddd,
		Number:   number,
		Type:     personentity.ContactTypeClient,
	}

	return &personentity.Contact{
		Entity:                  entity.NewEntity(),
		ContactCommonAttributes: contactCommonAttributes,
	}, nil
}
