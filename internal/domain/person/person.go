package personentity

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Person struct {
	entity.Entity
	PersonCommonAttributes
}

type PersonCommonAttributes struct {
	Name     string
	Email    string
	Cpf      string
	Birthday *time.Time
	Contact  *Contact
	Address  *addressentity.Address
}

func NewPerson(personCommonAttributes *PersonCommonAttributes) *Person {
	return &Person{
		Entity:                 entity.NewEntity(),
		PersonCommonAttributes: *personCommonAttributes,
	}
}

func (p *Person) AddContact(contact *Contact) error {
	p.Contact = contact
	p.Contact.ObjectID = p.ID
	return nil
}

func (p *Person) AddAddress(patchAddress *addressentity.Address) error {
	p.Address = patchAddress
	p.Address.ObjectID = p.ID
	return nil
}
