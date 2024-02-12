package userentity

import (
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
	Email    string          `bun:"column:email,unique,notnull" json:"email"`
	Password string          `bun:"-" json:"password"`
	Hash     string          `bun:"column:hash,notnull" json:"hash"`
	Schemas  []SchemaCompany `bun:"-" json:"schemas"`
}

type SchemaCompany struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
	Cnpj   string `json:"cnpj"`
}

func NewUser(userCommonAttributes UserCommonAttributes) *User {
	return &User{
		Entity:               entity.NewEntity(),
		UserCommonAttributes: userCommonAttributes,
	}
}

func (u *User) AddSchema(schema SchemaCompany) {
	u.Schemas = append(u.Schemas, schema)
}

func (u *User) RemoveSchema(schema SchemaCompany) {
	for i, s := range u.Schemas {
		if s == schema {
			u.Schemas = append(u.Schemas[:i], u.Schemas[i+1:]...)
			break
		}
	}
}
