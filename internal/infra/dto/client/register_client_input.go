package clientdto

import (
	"errors"
	"strings"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrNameRequired        = errors.New("name is required")
	ErrAddressRequired     = errors.New("address is required")
	ErrContactRequired     = errors.New("contact is required")
	ErrDeliveryTaxRequired = errors.New("delivery tax is required")
)

type CreateClientInput struct {
	personentity.PatchPerson
}

func (r *CreateClientInput) validate() error {
	if r.Name == nil || *r.Name == "" {
		return ErrNameRequired
	}
	if r.Contact == nil {
		return ErrContactRequired
	}
	if r.Address == nil {
		return ErrAddressRequired
	}

	if r.Address != nil && r.Address.DeliveryTax != nil && *r.Address.DeliveryTax == 0 {
		return ErrDeliveryTaxRequired
	}

	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *CreateClientInput) ToModel() (*cliententity.Client, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
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

	if err := person.AddContact(r.Contact, personentity.ContactTypeClient); err != nil {
		return nil, err
	}

	if err := person.AddAddress(r.Address); err != nil {
		return nil, err
	}

	return &cliententity.Client{
		Person: *person,
	}, nil
}
