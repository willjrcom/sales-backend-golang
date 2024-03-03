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
	Name     *string                `json:"name"`
	Email    *string                `json:"email"`
	Cpf      *string                `json:"cpf"`
	Birthday *time.Time             `json:"birthday"`
	Contact  *string                `json:"contact"`
	Address  *addressentity.Address `json:"address"`
}

func NewPerson(personCommonAttributes PersonCommonAttributes) *Person {
	return &Person{
		Entity:                 entity.NewEntity(),
		PersonCommonAttributes: personCommonAttributes,
	}
}

func (p *Person) AddContact(contactInput *string, contactType ContactType) error {
	ddd, number, err := ValidateAndExtractContact(*contactInput)

	if err != nil {
		return err
	}

	attributes := ContactCommonAttributes{
		ObjectID: p.ID,
		Ddd:      ddd,
		Number:   number,
		Type:     contactType,
	}

	p.Contact = NewContact(attributes)
	return nil
}

func (p *Person) AddAddress(addressCommonAttributes *addressentity.AddressCommonAttributes) error {
	addressCommonAttributes.ObjectID = p.ID
	p.Address = addressentity.NewAddress(addressCommonAttributes)

	if err := p.Address.Validate(); err != nil {
		return err
	}

	return nil
}
