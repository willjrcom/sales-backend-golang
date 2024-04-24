package employeedto

import (
	"errors"
	"strings"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UpdateEmployeeInput struct {
	personentity.PatchPerson
}

func (r *UpdateEmployeeInput) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *UpdateEmployeeInput) UpdateModel(client *employeeentity.Employee) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		client.Name = *r.Name
	}
	if r.Email != nil {
		client.Email = *r.Email
	}
	if r.Cpf != nil {
		client.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		client.Birthday = r.Birthday
	}
	if r.Contact != nil {
		if err := client.AddContact(r.Contact, personentity.ContactTypeEmployee); err != nil {
			return err
		}
	}
	if r.Address != nil {
		if err := client.AddAddress(&r.Address.AddressCommonAttributes); err != nil {
			return err
		}
	}

	return nil
}
