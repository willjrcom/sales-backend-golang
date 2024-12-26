package userdto

import (
	"errors"
	"strings"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UpdateUser struct {
	personentity.PatchPerson
}

func (r *UpdateUser) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *UpdateUser) UpdateModel(user *companyentity.User) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		user.Person.Name = *r.Name
	}
	if r.Email != nil {
		user.Email = *r.Email
	}
	if r.Cpf != nil {
		user.Person.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		user.Person.Birthday = r.Birthday
	}
	if r.Contact != nil {
		if err := user.AddContact(r.Contact, personentity.ContactTypeEmployee); err != nil {
			return err
		}
	} else {
		user.Contact = nil
	}

	if r.Address != nil {
		if err := user.AddAddress(r.Address); err != nil {
			return err
		}
	} else {
		user.Address = nil
	}

	return nil
}
