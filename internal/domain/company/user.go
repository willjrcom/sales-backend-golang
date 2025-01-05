package companyentity

import (
	"context"

	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserValue string

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
}

type UserCommonAttributes struct {
	personentity.Person
	Password       string             `bun:"-" json:"-"`
	Hash           string             `bun:"column:hash,notnull" json:"hash"`
	CompanyToUsers []CompanyToUsers   `bun:"rel:has-many,join:id=user_id" json:"company_users,omitempty"`
	Companies      []CompanyWithUsers `bun:"-" json:"companies"`
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
