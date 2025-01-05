package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Company struct {
	entity.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type CompanyCommonAttributes struct {
	SchemaName   string                 `bun:"schema_name,notnull"`
	BusinessName string                 `bun:"business_name,notnull"`
	TradeName    string                 `bun:"trade_name,notnull"`
	Cnpj         string                 `bun:"cnpj,notnull"`
	Email        string                 `bun:"email"`
	Contacts     []string               `bun:"contacts,type:jsonb"`
	Address      *addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull"`
}

type CompanyWithUsers struct {
	entity.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
	Users []User `bun:"m2m:company_to_users,join:CompanyWithUsers=User"`
}

type CompanyToUsers struct {
	CompanyWithUsersID uuid.UUID         `bun:"type:uuid,pk"`
	CompanyWithUsers   *CompanyWithUsers `bun:"rel:belongs-to,join:company_with_users_id=id"`
	UserID             uuid.UUID         `bun:"type:uuid,pk"`
	User               *User             `bun:"rel:belongs-to,join:user_id=id"`
}
