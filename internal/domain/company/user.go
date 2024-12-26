package companyentity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserValue string

type User struct {
	entity.Entity
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type UserCommonAttributes struct {
	PersonID       uuid.UUID           `bun:"column:person_id,type:uuid,notnull" json:"person_id"`
	Person         personentity.Person `bun:"rel:belongs-to,join:person_id=id" json:"person"`
	Email          string              `bun:"column:email,unique,notnull" json:"email"`
	Password       string              `bun:"-" json:"-"`
	Hash           string              `bun:"column:hash,notnull" json:"hash"`
	CompanyToUsers []CompanyToUsers    `bun:"rel:has-many,join:id=user_id" json:"company_users,omitempty"`
	Companies      []CompanyWithUsers  `bun:"-" json:"companies"`
}

func NewUser(userCommonAttributes UserCommonAttributes) *User {
	return &User{
		Entity:               entity.NewEntity(),
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
