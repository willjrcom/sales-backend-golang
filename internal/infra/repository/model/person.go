package model

import (
	"time"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Person struct {
	Name     string     `bun:"name,notnull"`
	Email    string     `bun:"email"`
	Cpf      string     `bun:"cpf"`
	Birthday *time.Time `bun:"birthday"`
	Contact  *Contact   `bun:"rel:has-one,join:id=object_id,notnull"`
	Address  *Address   `bun:"rel:has-one,join:id=object_id,notnull"`
}

func (p *Person) FromDomain(person *personentity.Person) {
	*p = Person{
		Name:     person.Name,
		Email:    person.Email,
		Cpf:      person.Cpf,
		Birthday: person.Birthday,
	}

	p.Contact.FromDomain(person.Contact)
	p.Address.FromDomain(person.Address)
}

func (p *Person) ToDomain() *personentity.Person {
	return &personentity.Person{
		PersonCommonAttributes: personentity.PersonCommonAttributes{
			Name:     p.Name,
			Email:    p.Email,
			Cpf:      p.Cpf,
			Birthday: p.Birthday,
			Contact:  p.Contact.ToDomain(),
			Address:  p.Address.ToDomain(),
		},
	}
}
