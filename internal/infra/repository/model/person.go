package model

import (
	"time"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Person struct {
	ImagePath string     `bun:"image_path"`
	Name      string     `bun:"name,notnull"`
	Email     string     `bun:"email"`
	Cpf       string     `bun:"cpf"`
	Birthday  *time.Time `bun:"birthday"`
	IsActive  bool       `bun:"column:is_active,type:boolean"`
	Contact   *Contact   `bun:"rel:has-one,join:id=object_id,notnull"`
	Address   *Address   `bun:"rel:has-one,join:id=object_id,notnull"`
}

func (p *Person) FromDomain(person *personentity.Person) {
	if person == nil {
		return
	}
	*p = Person{
		ImagePath: person.ImagePath,
		Name:      person.Name,
		Email:     person.Email,
		Cpf:       person.Cpf,
		Birthday:  person.Birthday,
		IsActive:  person.IsActive,
		Contact:   &Contact{},
		Address:   &Address{},
	}

	p.Contact.FromDomain(person.Contact)
	p.Address.FromDomain(person.Address)
}

func (p *Person) ToDomain() *personentity.Person {
	if p == nil {
		return nil
	}
	return &personentity.Person{
		PersonCommonAttributes: personentity.PersonCommonAttributes{
			ImagePath: p.ImagePath,
			Name:      p.Name,
			Email:     p.Email,
			Cpf:       p.Cpf,
			Birthday:  p.Birthday,
			IsActive:  p.IsActive,
			Contact:   p.Contact.ToDomain(),
			Address:   p.Address.ToDomain(),
		},
	}
}
