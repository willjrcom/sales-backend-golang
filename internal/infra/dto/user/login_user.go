package userdto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type LoginUserInput struct {
	companyentity.UserCommonAttributes
}

func (u *LoginUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if !utils.IsValidPassword(u.Password) && u.Password != "12345" {
		return ErrPasswordInvalid
	}

	return nil
}

func (u *LoginUserInput) ToModel() (*companyentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	userCommonAttributes := companyentity.UserCommonAttributes{
		Email:    u.Email,
		Password: u.Password,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
