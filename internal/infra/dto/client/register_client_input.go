package clientdto

import (
	"errors"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrAddressRequired  = errors.New("address is required")
	ErrContactsRequired = errors.New("contacts is required")
)

type RegisterClientInput struct {
	person   *personentity.Person
	Name     string     `json:"name"`
	Contact1 string     `json:"contact1"`
	Contact2 string     `json:"contact2"`
	Contact3 string     `json:"contact3"`
	Email    *string    `json:"email"`
	Cpf      *string    `json:"cpf"`
	Birthday *time.Time `json:"birthday"`
}

func (r *RegisterClientInput) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Contact1 == "" {
		return ErrContactsRequired
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
		r.person.Birthday = *r.Birthday
	}

	// Contacts
	if err := r.addContactToPerson(r.Contact1); err != nil {
		return nil, err
	}
	if err := r.addContactToPerson(r.Contact2); err != nil && r.Contact2 != "" {
		return nil, err
	}
	if err := r.addContactToPerson(r.Contact3); err != nil && r.Contact3 != "" {
		return nil, err
	}

	return &cliententity.Client{
		Person:      *r.person,
		TotalOrders: 0,
	}, nil
}

func (r *RegisterClientInput) addContactToPerson(c string) error {
	// Validate contact
	ddd, number, err := personentity.ValidateAndExtractContact(c)

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
	return nil
}
