package userdto

import (
	"errors"
	"strings"
	"time"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UserUpdateDTO struct {
	Name     *string                      `json:"name"`
	Email    *string                      `json:"email"`
	Cpf      *string                      `json:"cpf"`
	Birthday *time.Time                   `json:"birthday"`
	Contact  *contactdto.ContactUpdateDTO `json:"contact"`
	Address  *addressdto.AddressUpdateDTO `json:"address"`
}

func (r *UserUpdateDTO) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *UserUpdateDTO) UpdateDomain(user *companyentity.User) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		user.Person.Name = *r.Name
	}
	if r.Cpf != nil {
		user.Person.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		user.Person.Birthday = r.Birthday
	}
	if r.Contact != nil {
		r.Contact.UpdateDomain(user.Person.Contact)
	} else {
		user.Contact = nil
	}

	if r.Address != nil {
		r.Address.UpdateDomain(user.Address)
	} else {
		user.Address = nil
	}

	return nil
}
