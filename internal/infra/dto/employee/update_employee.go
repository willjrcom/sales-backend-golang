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

func (r *UpdateEmployeeInput) UpdateModel(employee *employeeentity.Employee) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		employee.User.Person.Name = *r.Name
	}
	if r.Email != nil {
		employee.User.Email = *r.Email
	}
	if r.Cpf != nil {
		employee.User.Person.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		employee.User.Person.Birthday = r.Birthday
	}
	if r.Contact != nil {
		if err := employee.User.Person.AddContact(r.Contact, personentity.ContactTypeEmployee); err != nil {
			return err
		}
	}
	if r.Address != nil {
		if err := employee.User.Person.AddAddress(r.Address); err != nil {
			return err
		}
	}

	return nil
}
