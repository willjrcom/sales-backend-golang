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
	Name      string                  `bun:"name,notnull" json:"name"`
	Email     string                  `bun:"email" json:"email,omitempty"`
	Cpf       string                  `bun:"cpf" json:"cpf,omitempty"`
	Birthday  *time.Time              `bun:"birthday" json:"birthday,omitempty"`
	Contacts  []Contact               `bun:"rel:has-many,join:id=object_id,notnull" json:"contacts,omitempty"`
	Addresses []addressentity.Address `bun:"rel:has-many,join:id=object_id,notnull" json:"addresses,omitempty"`
}

type PatchPerson struct {
	Name      *string                 `json:"name"`
	Email     *string                 `json:"email"`
	Cpf       *string                 `json:"cpf"`
	Birthday  *time.Time              `json:"birthday"`
	Contacts  []string                `json:"contacts"`
	Addresses []addressentity.Address `json:"addresses"`
}

func (p *Person) AddContact(contactInput string, contactType ContactType) error {
	ddd, number, err := ValidateAndExtractContact(contactInput)

	if err != nil {
		return err
	}

	contact := Contact{
		Entity: entity.NewEntity(),
		ContactCommonAttributes: ContactCommonAttributes{
			ObjectID: p.ID,
			Ddd:      ddd,
			Number:   number,
			Type:     contactType,
		},
	}

	p.Contacts = append(p.Contacts, contact)
	return nil
}
