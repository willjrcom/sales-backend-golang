package companyentity

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type UserValue string
type User struct {
	entity.Entity
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
}

type UserCommonAttributes struct {
	Email          string             `bun:"column:email,unique,notnull" json:"email"`
	Password       string             `bun:"-" json:"password"`
	Hash           string             `bun:"column:hash,notnull" json:"hash"`
	CompanyToUsers []CompanyToUsers   `bun:"rel:has-many,join:id=user_id" json:"company_users,omitempty"`
	Companies      []CompanyWithUsers `bun:"-" json:"companies"`
}

func NewUser(userCommonAttributes UserCommonAttributes) *User {
	return &User{
		Entity:               entity.NewEntity(),
		UserCommonAttributes: userCommonAttributes,
	}
}

func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	fmt.Println("Before updating:", u.Email)
	return nil
}

func (u *User) GetSchemas() []string {
	schemas := []string{}

	for _, company := range u.Companies {
		schemas = append(schemas, company.SchemaName)
	}

	return schemas
}
