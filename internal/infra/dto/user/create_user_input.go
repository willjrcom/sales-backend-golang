package userdto

import (
	"errors"
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrEmailInvalid             = errors.New("email is invalid")
	ErrMustHaveAtLeastOneSchema = errors.New("must have at least one schema")
)

type CreateUserInput struct {
	Email            string                      `json:"email"`
	Password         string                      `json:"password"`
	GeneratePassword bool                        `json:"generate_password"`
	Name             string                      `json:"name"`
	Cpf              string                      `json:"cpf,omitempty"`
	Birthday         *time.Time                  `json:"birthday,omitempty"`
	Contact          *personentity.Contact       `json:"contact,omitempty"`
	Address          *addressentity.PatchAddress `json:"address,omitempty"`
}

func (u *CreateUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if err := utils.ValidatePassword(u.Password); err != nil && !u.GeneratePassword {
		return err
	}

	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Cpf == "" {
		return errors.New("cpf is required")
	}
	if u.Birthday == nil {
		return errors.New("birthday is required")
	}
	if u.Contact == nil {
		return errors.New("contact is required")
	}
	if u.Address == nil {
		return errors.New("address is required")
	}
	return nil
}

func (u *CreateUserInput) ToModel() (*companyentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	if u.GeneratePassword {
		u.Password = utils.GeneratePassword(10, true, true, true)
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Name:     u.Name,
		Email:    u.Email,
		Cpf:      u.Cpf,
		Birthday: u.Birthday,
	}

	person := personentity.NewPerson(personCommonAttributes)

	if u.Contact != nil {
		if err := person.AddContact(&u.Contact.ContactCommonAttributes, personentity.ContactTypeEmployee); err != nil {
			return nil, err
		}
	}
	if u.Address != nil {
		if err := person.AddAddress(u.Address); err != nil {
			return nil, err
		}
	}

	return &companyentity.User{
		UserCommonAttributes: companyentity.UserCommonAttributes{
			Person:   *person,
			Password: u.Password,
		},
	}, nil
}
