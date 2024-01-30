package userentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type User struct {
	entity.Entity
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
}

type UserCommonAttributes struct {
	Email         string     `bun:"column:email,unique,notnull" json:"email"`
	Password      string     `bun:"-" json:"password"`
	Hash          string     `bun:"column:hash,notnull" json:"hash"`
	Schemas       []string   `bun:"column:schemas,notnull,type:jsonb,notnull" json:"schemas"`
	CurrentSchema *string    `bun:"-" json:"current_schema"`
	EmployeeID    *uuid.UUID `bun:"column:employee_id,type:uuid,notnull" json:"employee_id"`
}

func NewUser(userCommonAttributes UserCommonAttributes) *User {
	return &User{
		Entity:               entity.NewEntity(),
		UserCommonAttributes: userCommonAttributes,
	}
}

func (u *User) AddSchema(schema string) {
	u.Schemas = append(u.Schemas, schema)
}

func (u *User) RemoveSchema(schema string) {
	for i, s := range u.Schemas {
		if s == schema {
			u.Schemas = append(u.Schemas[:i], u.Schemas[i+1:]...)
			break
		}
	}
}
