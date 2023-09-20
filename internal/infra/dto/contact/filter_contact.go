package contactdto

import (
	"errors"

	"github.com/google/uuid"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidFilter = errors.New("invalid filter")
)

type FilterContact struct {
	PersonID *uuid.UUID `json:"person_id"`
	Contact  *string    `json:"contact"`
}

func (f *FilterContact) validate() error {
	if f.PersonID == nil && f.Contact == nil {
		return ErrInvalidFilter
	}

	return nil
}
func (f *FilterContact) ToModel() (*personentity.Contact, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	contact := &personentity.Contact{}

	if f.PersonID != nil {
		contact.PersonID = *f.PersonID
	}

	if f.Contact == nil {
		return contact, nil
	}

	ddd, number, err := personentity.ValidateAndExtractContact(*f.Contact)

	if err != nil {
		return nil, err
	}

	contact.Ddd = ddd
	contact.Number = number
	return contact, nil
}
