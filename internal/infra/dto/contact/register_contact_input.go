package contactdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrContactIsEmpty = errors.New("contact is required")
	ErrIDIsEmpty      = errors.New("ID is required")
)

type RegisterContactInput struct {
	ClientID   uuid.UUID `json:"client_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	Contact    string    `json:"contact"`
}

func (c *RegisterContactInput) validate() error {
	if c.Contact == "" {
		return ErrContactIsEmpty
	}
	if c.ClientID == uuid.Nil && c.EmployeeID == uuid.Nil {
		return ErrIDIsEmpty
	}

	return nil
}

func (c *RegisterContactInput) ToModel() (*personentity.Contact, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	// Validate contact
	ddd, number, err := personentity.ValidateAndExtractContact(c.Contact)

	if err != nil {
		return nil, err
	}

	id := uuid.Nil
	contactType := personentity.ContactTypeClient

	if c.ClientID != uuid.Nil {
		id = c.ClientID
	} else if c.EmployeeID != uuid.Nil {
		id = c.EmployeeID
		contactType = personentity.ContactTypeEmployee
	}

	if id == uuid.Nil {
		return nil, ErrIDIsEmpty
	}

	contactCommonAttributes := personentity.ContactCommonAttributes{
		ObjectID: id,
		Ddd:      ddd,
		Number:   number,
		Type:     contactType,
	}

	return &personentity.Contact{
		Entity:                  entity.NewEntity(),
		ContactCommonAttributes: contactCommonAttributes,
	}, nil
}
