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
	Password       string             `bun:"-" json:"-"`
	Hash           string             `bun:"column:hash,notnull" json:"hash"`
	CompanyToUsers []CompanyToUsers   `bun:"rel:has-many,join:id=user_id" json:"company_users,omitempty"`
	Companies      []CompanyWithUsers `bun:"-" json:"companies"`
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
			},
			Password:       user.Password,
			CompanyToUsers: []CompanyToUsers{},
			Companies:      []CompanyWithUsers{},
		},
	}

	for _, company := range user.Companies {
		c := CompanyWithUsers{}
		c.FromDomain(&company)
		u.Companies = append(u.Companies, c)
	}

	for _, companyToUser := range user.CompanyToUsers {
		c := CompanyToUsers{}
		c.FromDomain(&companyToUser)
		u.CompanyToUsers = append(u.CompanyToUsers, c)
	}
}

func (u *User) ToDomain() *companyentity.User {
	return &companyentity.User{
		Entity: u.Entity.ToDomain(),
		UserCommonAttributes: companyentity.UserCommonAttributes{
			Person:         *u.Person.ToDomain(),
			Password:       u.Password,
			Hash:           u.Hash,
			CompanyToUsers: []companyentity.CompanyToUsers{},
			Companies:      []companyentity.CompanyWithUsers{},
		},
	}
}
