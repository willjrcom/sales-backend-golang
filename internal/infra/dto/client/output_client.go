package clientdto

import (
	"time"

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
	Name     string                     `json:"name"`
	Birthday time.Time                  `json:"birthday"`
	Email    string                     `json:"email"`
	Cpf      string                     `json:"cpf"`
	Contacts []contactdto.ContactOutput `json:"contacts,omitempty"`
	Address  []addressdto.AddressOutput `json:"address,omitempty"`
}

func (p *Person) FromModel(model *personentity.Person) {
	p.Name = model.Name
	p.Birthday = *model.Birthday
	p.Email = model.Email
	p.Cpf = model.Cpf

}

func (c *ClientOutput) FromModel(model *cliententity.Client) {
	c.Person.FromModel(&model.Person)
	c.TotalOrders = model.TotalOrders
}
