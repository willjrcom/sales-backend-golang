package clientdto

import (
	"errors"
	"strings"
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type ClientUpdateDTO struct {
	Name     *string                      `json:"name"`
	Email    *string                      `json:"email"`
	Cpf      *string                      `json:"cpf"`
	Birthday *time.Time                   `json:"birthday"`
	Contact  *contactdto.ContactUpdateDTO `json:"contact"`
	Address  *addressdto.AddressUpdateDTO `json:"address"`
}

func (r *ClientUpdateDTO) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	if r.Address != nil && r.Address.DeliveryTax != nil && r.Address.DeliveryTax.IsZero() {
		return ErrDeliveryTaxRequired
	}

	return nil
}

func (r *ClientUpdateDTO) UpdateDomain(client *cliententity.Client) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		client.Name = *r.Name
	}
	if r.Email != nil {
		client.Email = *r.Email
	}
	if r.Cpf != nil {
		client.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		client.Birthday = r.Birthday
	}
	if r.Contact != nil {
		if client.Contact == nil {
			client.Contact = &personentity.Contact{
				Entity:   entity.NewEntity(),
				ObjectID: client.ID,
			}
		}
		r.Contact.UpdateDomain(client.Contact, personentity.ContactTypeClient)
	}
	if r.Address != nil {
		if client.Address == nil {
			client.Address = &addressentity.Address{
				Entity:   entity.NewEntity(),
				ObjectID: client.ID,
			}
		}
		r.Address.UpdateDomain(client.Address)
	}

	return nil
}
