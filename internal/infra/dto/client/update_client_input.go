package clientdto

import (
	"errors"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
)

type UpdateClientInput struct {
	Name     *string    `json:"name"`
	Email    *string    `json:"email"`
	Cpf      *string    `json:"cpf"`
	Birthday *time.Time `json:"birthday"`
}

func (r *UpdateClientInput) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return errors.New("name is required")
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
	if r.Email != nil {
		client.Email = *r.Email
	}
	if r.Cpf != nil {
		client.Cpf = *r.Cpf
	}
	if r.Birthday != nil {
		client.Birthday = *r.Birthday
	}

	return nil
}
