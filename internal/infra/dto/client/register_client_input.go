package clientdto

import (
	"errors"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrAddressRequired  = errors.New("address is required")
	ErrMaxAddresses     = errors.New("max addresses is 3")
	ErrContactsRequired = errors.New("contacts is required")
	ErrMaxContacts      = errors.New("max contacts is 3")
)

type RegisterClientInput struct {
	person    *personentity.Person
	Name      string                            `json:"name"`
	Email     *string                           `json:"email"`
	Cpf       *string                           `json:"cpf"`
	Birthday  *time.Time                        `json:"birthday"`
	Contacts  []string                          `json:"contacts"`
	Addresses []addressdto.RegisterAddressInput `json:"addresses"`
}

func (r *RegisterClientInput) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if len(r.Contacts) == 0 {
		return ErrContactsRequired
	}
	if len(r.Contacts) > 3 {
		return ErrMaxContacts
	}
	if len(r.Addresses) == 0 {
		return ErrAddressRequired
	}
	if len(r.Addresses) > 3 {
		return ErrMaxAddresses
	}
	if len(r.Addresses) == 1 {
		r.Addresses[0].IsDefault = true
	}

	return nil
}

func (r *RegisterClientInput) ToModel() (*cliententity.Client, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	// Create person
	r.person = &personentity.Person{
		Entity: entity.NewEntity(),
		Name:   r.Name,
	}

	// Optional fields
	if r.Email != nil {
		r.person.Email = *r.Email
	}
	if r.Cpf != nil {
		r.person.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		r.person.Birthday = r.Birthday
	}

	// Contacts
	if err := r.addContactsToPerson(r.Contacts); err != nil {
		return nil, err
	}

	if err := r.addAddressesToPerson(r.Addresses); err != nil {
		return nil, err
	}

	return &cliententity.Client{
		Person:      *r.person,
		TotalOrders: 0,
	}, nil
}

func (r *RegisterClientInput) addContactsToPerson(list []string) error {
	for _, contact := range list {
		// Validate contact
		ddd, number, err := personentity.ValidateAndExtractContact(contact)

		if err != nil {
			return err
		}

		// Create first contact
		contact := &personentity.Contact{
			Entity:   entity.NewEntity(),
			Ddd:      ddd,
			Number:   number,
			PersonID: r.person.ID,
		}

		r.person.Contacts = append(r.person.Contacts, *contact)
	}
	return nil
}

func (r *RegisterClientInput) addAddressesToPerson(list []addressdto.RegisterAddressInput) error {
	for _, dto := range list {
		address, err := dto.ToModel()

		if err != nil {
			return err
		}

		address.PersonID = r.person.ID
		r.person.Addresses = append(r.person.Addresses, *address)
	}

	return nil
}
