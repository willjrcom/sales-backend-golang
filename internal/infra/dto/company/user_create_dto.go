package companydto

import (
	"errors"
	"time"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

var (
	ErrMustBeName     = errors.New("name is required")
	ErrMustBeEmail    = errors.New("email is required")
	ErrMustBePassword = errors.New("password is required")
	ErrMustBeCpf      = errors.New("cpf is required")
	ErrMustBeBirthday = errors.New("birthday is required")
	ErrMustBeContact  = errors.New("contact is required")
	ErrMustBeAddress  = errors.New("address is required")
)

type UserCreateDTO struct {
	Name     string                       `json:"name"`
	Email    string                       `json:"email"`
	Cpf      string                       `json:"cpf"`
	Birthday *time.Time                   `json:"birthday"`
	Contact  *contactdto.ContactCreateDTO `json:"contact"`
	Address  *addressdto.AddressCreateDTO `json:"address"`
	Password string                       `json:"password"`
}

func (c *UserCreateDTO) validate() error {
	if c.Name == "" {
		return ErrMustBeName
	}
	if c.Email == "" {
		return ErrMustBeEmail
	}
	if c.Cpf == "" {
		return ErrMustBeCpf
	}
	if c.Birthday == nil {
		return ErrMustBeBirthday
	}
	if c.Contact == nil {
		return ErrMustBeContact
	}
	if c.Address == nil {
		return ErrMustBeAddress
	}
	if c.Password == "" {
		return ErrMustBePassword
	}
	return nil
}

func (c *UserCreateDTO) ToDomain() (*companyentity.User, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Name:     c.Name,
		Email:    c.Email,
		Cpf:      c.Cpf,
		Birthday: c.Birthday,
	}

	contact, err := c.Contact.ToModel()
	if err != nil {
		return nil, err
	}

	personCommonAttributes.Contact = contact

	address, err := c.Address.ToModel(false)
	if err != nil {
		return nil, err
	}

	personCommonAttributes.Address = address

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person: *person,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
