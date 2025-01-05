package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserValue string

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type UserCommonAttributes struct {
	personentity.Person
	Password       string             `bun:"-" json:"-"`
	Hash           string             `bun:"column:hash,notnull" json:"hash"`
	CompanyToUsers []CompanyToUsers   `bun:"rel:has-many,join:id=user_id" json:"company_users,omitempty"`
	Companies      []CompanyWithUsers `bun:"-" json:"companies"`
}

func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	return nil
}
