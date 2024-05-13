package employeedto

import (
	"errors"
	"strings"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrNameRequired    = errors.New("name is required")
	ErrAddressRequired = errors.New("address is required")
	ErrContactRequired = errors.New("contact is required")
)

type CreateEmployeeInput struct {
	personentity.PatchPerson
}

func (r *CreateEmployeeInput) validate() error {
	if r.Name == nil || *r.Name == "" {
		return ErrNameRequired
	}
	if r.Contact == nil {
		return ErrContactRequired
	}
	if r.Address == nil {
		return ErrAddressRequired
	}

	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *CreateEmployeeInput) ToModel() (*employeeentity.Employee, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := personentity.PersonCommonAttributes{
		Name: *r.Name,
	}

	// Create person
	person := personentity.NewPerson(personCommonAttributes)

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

	if err := person.AddContact(r.Contact, personentity.ContactTypeEmployee); err != nil {
		return nil, err
	}

	if err := person.AddAddress(&r.Address.AddressCommonAttributes); err != nil {
		return nil, err
	}

	return &employeeentity.Employee{
		Person: *person,
	}, nil
}
