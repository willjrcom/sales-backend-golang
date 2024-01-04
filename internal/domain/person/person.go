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
	Email     string                  `bun:"email" json:"email"`
	Cpf       string                  `bun:"cpf,unique" json:"cpf"`
	Birthday  *time.Time              `bun:"birthday" json:"birthday"`
	Contacts  []Contact               `bun:"rel:has-many,join:id=person_id,notnull" json:"contacts"`
	Addresses []addressentity.Address `bun:"rel:has-many,join:id=person_id,notnull" json:"addresses"`
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
		ContactCommonAttributes: ContactCommonAttributes{
			PersonID: p.ID,
			Ddd:      ddd,
			Number:   number,
			Type:     contactType,
		},
	}

	p.Contacts = append(p.Contacts, contact)
	return nil
}
