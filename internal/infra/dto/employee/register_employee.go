package employeedto

import (
	"errors"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrAddressRequired  = errors.New("address is required")
	ErrMaxAddresses     = errors.New("max addresses is 3")
	ErrContactsRequired = errors.New("contacts is required")
	ErrMaxContacts      = errors.New("max contacts is 3")
)

type RegisterEmployeeInput struct {
	personentity.PatchPerson
}

func (r *RegisterEmployeeInput) Validate() error {
	if r.Name == nil || *r.Name == "" {
		return ErrNameRequired
	}
	if len(r.Contacts) == 0 {
		return ErrContactsRequired
	}
	if len(r.Contacts) > 3 {
		return ErrMaxContacts
	}
	if len(r.Addresses) == 0 {
		return ErrAddressRequired
	}
	if len(r.Addresses) > 3 {
		return ErrMaxAddresses
	}
	if len(r.Addresses) == 1 {
		r.Addresses[0].IsDefault = true
	}

	return nil
}

func (r *RegisterEmployeeInput) ToModel() (*employeeentity.Employee, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := personentity.PersonCommonAttributes{
		Name:      *r.Name,
		Addresses: r.Addresses,
	}

	// Create person
	person := &personentity.Person{
		Entity:                 entity.NewEntity(),
		PersonCommonAttributes: personCommonAttributes,
	}

	// Optional fields
	if r.Email != nil {
		person.Email = *r.Email
	}
	if r.Cpf != nil {
		person.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		person.Birthday = r.Birthday
	}

	for _, contact := range r.Contacts {
		if err := person.AddContact(contact, personentity.ContactTypeClient); err != nil {
			return nil, err
		}
	}

	return &employeeentity.Employee{
		Person: *person,
	}, nil
}
