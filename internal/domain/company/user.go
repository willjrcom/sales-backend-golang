package companyentity

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserValue string

type User struct {
	entity.Entity
	UserCommonAttributes
}

type UserCommonAttributes struct {
	personentity.Person
	Password       string
	Hash           string
	CompanyToUsers []CompanyToUsers
	Companies      []CompanyWithUsers
}

func NewUser(userCommonAttributes *UserCommonAttributes) *User {
	return &User{
		Entity:               entity.NewEntity(),
		UserCommonAttributes: *userCommonAttributes,
	}
}

func (u *User) GetSchemas() []string {
	schemas := []string{}

	for _, company := range u.Companies {
		schemas = append(schemas, company.SchemaName)
	}

	return schemas
}

func (p *User) AddContact(contact *personentity.Contact) error {
	p.Contact = contact
	p.Contact.ObjectID = p.ID
	return nil
}

func (p *User) AddAddress(address *addressentity.Address) error {
	p.Address = address
	p.Address.ObjectID = p.ID
	return nil
}
