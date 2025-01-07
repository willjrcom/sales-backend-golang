package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CompanyToUsers struct {
	bun.BaseModel `bun:"table:company_to_users"`
	CompanyID     uuid.UUID `bun:"type:uuid,pk"`
	Company       *Company  `bun:"rel:belongs-to,join:company_id=id"`
	UserID        uuid.UUID `bun:"type:uuid,pk"`
	User          *User     `bun:"rel:belongs-to,join:user_id=id"`
}
