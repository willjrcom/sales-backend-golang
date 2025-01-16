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
	bun.BaseModel `bun:"table:public.users,alias:u"`
	UserCommonAttributes
}

type UserCommonAttributes struct {
	PublicPerson
	Password  string    `bun:"-"`
	Hash      string    `bun:"column:hash,notnull"`
	Companies []Company `bun:"m2m:public.company_to_users,join:User=Company"`
}

func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	return nil
}

func (u *User) FromDomain(user *companyentity.User) {
	if user == nil {
		return
	}
	*u = User{
		Entity: entitymodel.FromDomain(user.Entity),
		UserCommonAttributes: UserCommonAttributes{
			PublicPerson: PublicPerson{
				Name:     user.Name,
				Email:    user.Email,
				Cpf:      user.Cpf,
				Birthday: user.Birthday,
				Contact:  &PublicContact{},
				Address:  &PublicAddress{},
			},
			Password:  user.Password,
			Hash:      user.Hash,
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

	if user.Contact == nil {
		u.Contact = nil
	}
	if user.Address == nil {
		u.Address = nil
	}
}

func (u *User) ToDomain() *companyentity.User {
	if u == nil {
		return nil
	}
	user := &companyentity.User{
		Entity: u.Entity.ToDomain(),
		UserCommonAttributes: companyentity.UserCommonAttributes{
			Person:    *u.PublicPerson.ToDomain(),
			Password:  u.Password,
			Hash:      u.Hash,
			Companies: []companyentity.Company{},
		},
	}

	for _, company := range u.Companies {
		c := company.ToDomain()
		user.Companies = append(user.Companies, *c)
	}

	return user
}
