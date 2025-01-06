package model

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyToUsers struct {
	CompanyWithUsersID uuid.UUID         `bun:"type:uuid,pk"`
	CompanyWithUsers   *CompanyWithUsers `bun:"rel:belongs-to,join:company_with_users_id=id"`
	UserID             uuid.UUID         `bun:"type:uuid,pk"`
	User               *User             `bun:"rel:belongs-to,join:user_id=id"`
}

func (c *CompanyToUsers) FromDomain(model *companyentity.CompanyToUsers) {
	*c = CompanyToUsers{
		CompanyWithUsersID: model.CompanyWithUsersID,
		CompanyWithUsers:   nil,
		UserID:             model.UserID,
		User:               nil,
	}
}

func (c *CompanyToUsers) ToDomain() *companyentity.CompanyToUsers {
	if c == nil {
		return nil
	}
	return &companyentity.CompanyToUsers{
		CompanyWithUsersID: c.CompanyWithUsersID,
		CompanyWithUsers:   c.CompanyWithUsers.ToDomain(),
		UserID:             c.UserID,
		User:               c.User.ToDomain(),
	}
}
