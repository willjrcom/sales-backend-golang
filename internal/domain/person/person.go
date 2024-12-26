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
	Name     string                 `bun:"name,notnull" json:"name"`
	Email    string                 `bun:"email" json:"email,omitempty"`
	Cpf      string                 `bun:"cpf" json:"cpf,omitempty"`
	Birthday *time.Time             `bun:"birthday" json:"birthday,omitempty"`
	Contact  *Contact               `bun:"rel:has-one,join:id=object_id,notnull" json:"contact,omitempty"`
	Address  *addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull" json:"address,omitempty"`
}

type PatchPerson struct {
	Name     *string                     `json:"name"`
	Email    *string                     `json:"email"`
	Cpf      *string                     `json:"cpf"`
	Birthday *time.Time                  `json:"birthday"`
	Contact  *ContactCommonAttributes    `json:"contact"`
	Address  *addressentity.PatchAddress `json:"address"`
}

func NewPerson(personCommonAttributes *PersonCommonAttributes) *Person {
	return &Person{
		Entity:                 entity.NewEntity(),
		PersonCommonAttributes: *personCommonAttributes,
	}
}

func (p *Person) AddContact(contactInput *ContactCommonAttributes, contactType ContactType) error {
	contactInput.Type = contactType

	p.Contact = NewContact(contactInput)
	p.Contact.ObjectID = p.ID
	return nil
}

func (p *Person) AddAddress(patchAddress *addressentity.PatchAddress) error {
	p.Address = addressentity.NewAddressFromPatch(patchAddress, p.ID)

	if err := p.Address.Validate(); err != nil {
		return err
	}

	return nil
}
