package model

import (
	"time"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type PublicPerson struct {
	ImagePath string         `bun:"image_path"`
	Name      string         `bun:"name,notnull"`
	Email     string         `bun:"email"`
	Cpf       string         `bun:"cpf"`
	Birthday  *time.Time     `bun:"birthday"`
	Contact   *PublicContact `bun:"rel:has-one,join:id=object_id,notnull"`
	Address   *PublicAddress `bun:"rel:has-one,join:id=object_id,notnull"`
	IsActive  bool           `bun:"column:is_active,type:boolean"`
}

func (p *PublicPerson) FromDomain(person *personentity.Person) {
	if person == nil {
		return
	}
	*p = PublicPerson{
		ImagePath: person.ImagePath,
		Name:      person.Name,
		Email:     person.Email,
		Cpf:       person.Cpf,
		Birthday:  person.Birthday,
		Contact:   &PublicContact{},
		Address:   &PublicAddress{},
		IsActive:  p.IsActive,
	}

	p.Contact.FromDomain(person.Contact)
	p.Address.FromDomain(person.Address)
}

func (p *PublicPerson) ToDomain() *personentity.Person {
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
			Contact:   p.Contact.ToDomain(),
			Address:   p.Address.ToDomain(),
			IsActive:  p.IsActive,
		},
	}
}
