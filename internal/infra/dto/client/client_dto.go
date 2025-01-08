package clientdto

import (
	"time"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type ClientDTO struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Cpf      string                 `json:"cpf"`
	Birthday *time.Time             `json:"birthday"`
	Contact  *contactdto.ContactDTO `json:"contact"`
	Address  *addressdto.AddressDTO `json:"address"`
}

func (c *ClientDTO) FromDomain(client *cliententity.Client) {
	if client == nil {
		return
	}
	*c = ClientDTO{
		ID:       client.ID,
		Name:     client.Name,
		Email:    client.Email,
		Cpf:      client.Cpf,
		Birthday: client.Birthday,
		Contact:  &contactdto.ContactDTO{},
		Address:  &addressdto.AddressDTO{},
	}

	c.Contact.FromDomain(client.Contact)
	c.Address.FromDomain(client.Address)

	if client.Contact == nil {
		c.Contact = nil
	}
	if client.Address == nil {
		c.Address = nil
	}
}
