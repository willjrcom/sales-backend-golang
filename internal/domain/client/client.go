package cliententity

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	entity.Entity
	personentity.Person
}

func NewClient(person *personentity.Person) *Client {
	return &Client{
		Entity: entity.NewEntity(),
		Person: *person,
	}
}

func (p *Client) AddContact(contact *personentity.Contact) error {
	p.Contact = contact
	p.Contact.ObjectID = p.ID
	return nil
}

func (p *Client) AddAddress(address *addressentity.Address) error {
	p.Address = address
	p.Address.ObjectID = p.ID
	return nil
}
