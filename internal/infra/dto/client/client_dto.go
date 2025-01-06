package clientdto

import (
	"time"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type ClientOutput struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Cpf      string                 `json:"cpf"`
	Birthday *time.Time             `json:"birthday"`
	Contact  *contactdto.ContactDTO `json:"contact"`
	Address  *addressdto.AddressDTO `json:"address"`
}

func (c *ClientOutput) FromModel(model *cliententity.Client) {
	*c = ClientOutput{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Cpf:      model.Cpf,
		Birthday: model.Birthday,
	}

	c.Contact.FromDomain(model.Contact)
	c.Address.FromModel(model.Address)
}
