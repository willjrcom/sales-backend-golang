package clientdto

import (
	"errors"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
)

type UpdateClientInput struct {
	Name         *string    `json:"name"`
	Birthday     *time.Time `json:"birthday"`
	Email        *string    `json:"email"`
	Contacts     *[]string  `json:"contacts"`
	Cpf          *string    `json:"cpf"`
	Street       *string    `json:"street"`
	Number       *string    `json:"number"`
	Complement   *string    `json:"complement"`
	Reference    *string    `json:"reference"`
	Neighborhood *string    `json:"neighborhood"`
	City         *string    `json:"city"`
	State        *string    `json:"state"`
	Cep          *string    `json:"cep"`
}

func (r *UpdateClientInput) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return errors.New("name is required")
	}
	if r.Street != nil && *r.Street == "" {
		return errors.New("street is required")
	}
	if r.Number != nil && *r.Number == "" {
		return errors.New("number is required")
	}
	if r.Neighborhood != nil && *r.Neighborhood == "" {
		return errors.New("Address is required")
	}

	return nil
}

func (r *UpdateClientInput) UpdateModel(client *cliententity.Client) error {
	if err := r.Validate(); err != nil {
		return err
	}

	if r.Name != nil {
		client.Name = *r.Name
	}
	if r.Birthday != nil {
		client.Birthday = *r.Birthday
	}
	if r.Email != nil {
		client.Email = *r.Email
	}
	if r.Contacts != nil {
		client.Contacts = *r.Contacts
	}
	if r.Cpf != nil {
		client.Cpf = *r.Cpf
	}
	if r.Street != nil {
		client.Address.Street = *r.Street
	}
	if r.Number != nil {
		client.Address.Number = *r.Number
	}
	if r.Complement != nil {
		client.Address.Complement = *r.Complement
	}
	if r.Reference != nil {
		client.Address.Reference = *r.Reference
	}
	if r.Neighborhood != nil {
		client.Address.Neighborhood = *r.Neighborhood
	}
	if r.City != nil {
		client.Address.City = *r.City
	}
	if r.State != nil {
		client.Address.State = *r.State
	}
	if r.Cep != nil {
		client.Address.Cep = *r.Cep
	}

	return nil
}
