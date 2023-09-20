package clientdto

import (
	"time"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type ClientOutput struct {
	Person
	TotalOrders int `json:"total_orders"`
}

type Person struct {
	ID        uuid.UUID                  `json:"id"`
	Name      string                     `json:"name"`
	Email     string                     `json:"email"`
	Cpf       string                     `json:"cpf"`
	Birthday  *time.Time                 `json:"birthday,omitempty"`
	Contacts  []contactdto.ContactOutput `json:"contacts"`
	Addresses []addressdto.AddressOutput `json:"addresses"`
}

func (p *Person) FromModel(model *personentity.Person) {
	p.Name = model.Name

	if model.Birthday != nil {
		p.Birthday = model.Birthday
	}
	p.Email = model.Email
	p.Cpf = model.Cpf

	p.Contacts = make([]contactdto.ContactOutput, len(model.Contacts))
	p.Addresses = make([]addressdto.AddressOutput, len(model.Addresses))

	for i, v := range model.Contacts {
		p.Contacts[i].FromModel(&v)
	}

	for i, v := range model.Addresses {
		p.Addresses[i].FromModel(&v)
	}
}

func (c *ClientOutput) FromModel(model *cliententity.Client) {
	c.ID = model.ID
	c.Person.FromModel(&model.Person)
	c.TotalOrders = model.TotalOrders
}
