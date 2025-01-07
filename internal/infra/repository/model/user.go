package model

import (
	"context"

	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type UserValue string

type User struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
}

type UserCommonAttributes struct {
	Person
	Password  string    `bun:"-"`
	Hash      string    `bun:"column:hash,notnull"`
	Companies []Company `bun:"m2m:company_to_users,join:User=Company"`
}

func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	return nil
}

func (u *User) FromDomain(user *companyentity.User) {
	*u = User{
		Entity: entitymodel.FromDomain(user.Entity),
		UserCommonAttributes: UserCommonAttributes{
			Person: Person{
				Name:     user.Person.Name,
				Email:    user.Person.Email,
				Cpf:      user.Person.Cpf,
				Birthday: user.Person.Birthday,
				Contact:  &Contact{},
				Address:  &Address{},
			},
			Password:  user.Password,
			Companies: []Company{},
		},
	}

	for _, company := range user.Companies {
		c := Company{}
		c.FromDomain(&company)
		u.Companies = append(u.Companies, c)
	}

	u.Contact.FromDomain(user.Contact)
	u.Address.FromDomain(user.Address)
}

func (u *User) ToDomain() *companyentity.User {
	if u == nil {
		return nil
	}
	user := &companyentity.User{
		Entity: u.Entity.ToDomain(),
		UserCommonAttributes: companyentity.UserCommonAttributes{
			Person:    *u.Person.ToDomain(),
			Password:  u.Password,
			Hash:      u.Hash,
			Companies: []companyentity.Company{},
		},
	}

	for _, company := range u.Companies {
		c := company.ToDomain()
		user.Companies = append(user.Companies, *c)
	}

	for _, companyToUser := range u.Companies {
		c := companyToUser.ToDomain()
		user.Companies = append(user.Companies, *c)
	}

	return user
}
