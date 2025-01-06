package clientdto

import (
	"errors"
	"strings"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

var (
	ErrNameRequired        = errors.New("name is required")
	ErrAddressRequired     = errors.New("address is required")
	ErrContactRequired     = errors.New("contact is required")
	ErrDeliveryTaxRequired = errors.New("delivery tax is required")
)

type ClientCreateDTO struct {
	Name     string                       `json:"name"`
	Email    *string                      `json:"email"`
	Cpf      *string                      `json:"cpf"`
	Birthday *time.Time                   `json:"birthday"`
	Contact  *contactdto.ContactCreateDTO `json:"contact"`
	Address  *addressdto.AddressCreateDTO `json:"address"`
}

func (r *ClientCreateDTO) validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}
	if r.Contact == nil {
		return ErrContactRequired
	}
	if r.Address == nil {
		return ErrAddressRequired
	}

	return nil
}

func (r *ClientCreateDTO) ToDomain() (*cliententity.Client, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Name: r.Name,
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

	// Contact
	contact, err := r.Contact.ToDomain()
	if err != nil {
		return nil, err
	}

	contact.Type = personentity.ContactTypeClient

	if err := person.AddContact(contact); err != nil {
		return nil, err
	}

	// Address
	address, err := r.Address.ToDomain(true)
	if err != nil {
		return nil, err
	}

	if err := person.AddAddress(address); err != nil {
		return nil, err
	}

	return cliententity.NewClient(person), nil
}
