package clientdto

import (
	"errors"
	"strings"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UpdateClientInput struct {
	personentity.PatchPerson
}

func (r *UpdateClientInput) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	if r.Address != nil && r.Address.DeliveryTax != nil && *r.Address.DeliveryTax == 0 {
		return ErrDeliveryTaxRequired
	}

	return nil
}

func (r *UpdateClientInput) UpdateModel(client *cliententity.Client) error {
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
		if err := client.AddContact(r.Contact, personentity.ContactTypeEmployee); err != nil {
			return err
		}
	}
	if r.Address != nil {
		if err := client.AddAddress(r.Address); err != nil {
			return err
		}
	}

	return nil
}
