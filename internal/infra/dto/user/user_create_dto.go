package userdto

import (
	"errors"
	"time"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrEmailInvalid             = errors.New("email is invalid")
	ErrMustHaveAtLeastOneSchema = errors.New("must have at least one schema")
)

type UserCreateDTO struct {
	Email            string                       `json:"email"`
	Password         string                       `json:"password"`
	GeneratePassword bool                         `json:"generate_password"`
	Name             string                       `json:"name"`
	Cpf              string                       `json:"cpf,omitempty"`
	Birthday         *time.Time                   `json:"birthday,omitempty"`
	Contact          *contactdto.ContactCreateDTO `json:"contact,omitempty"`
	Address          *addressdto.AddressCreateDTO `json:"address,omitempty"`
}

func (u *UserCreateDTO) validate() error {
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

func (u *UserCreateDTO) ToModel() (*companyentity.User, error) {
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
		contact, err := u.Contact.ToModel()
		if err != nil {
			return nil, err
		}
		if err := person.AddContact(contact); err != nil {
			return nil, err
		}
	}
	if u.Address != nil {
		address, err := u.Address.ToModel(false)
		if err != nil {
			return nil, err
		}

		if err := person.AddAddress(address); err != nil {
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
