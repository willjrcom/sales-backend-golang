package companyentity

import (
	"context"

	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserValue string

type User struct {
	UserCommonAttributes
}

type UserCommonAttributes struct {
	personentity.Person
	Password       string
	Hash           string
	CompanyToUsers []CompanyToUsers
	Companies      []CompanyWithUsers
}

func NewUser(userCommonAttributes UserCommonAttributes) *User {
	return &User{
		UserCommonAttributes: userCommonAttributes,
	}
}

func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	return nil
}

func (u *User) GetSchemas() []string {
	schemas := []string{}

	for _, company := range u.Companies {
		schemas = append(schemas, company.SchemaName)
	}

	return schemas
}
