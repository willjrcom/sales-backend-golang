package model

import (
	"context"

	"github.com/uptrace/bun"
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
