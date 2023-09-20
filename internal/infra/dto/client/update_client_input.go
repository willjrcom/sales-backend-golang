package clientdto

import (
	"errors"
	"strings"
	"time"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UpdateClientInput struct {
	Name     *string    `json:"name"`
	Email    *string    `json:"email"`
	Cpf      *string    `json:"cpf"`
	Birthday *time.Time `json:"birthday"`
}

func (r *UpdateClientInput) Validate() error {
	if r.Email != nil && strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
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
		client.Birthday = r.Birthday
	}

	return nil
}
