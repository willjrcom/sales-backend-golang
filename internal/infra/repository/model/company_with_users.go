package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyWithUsers struct {
	entitymodel.Entity
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
