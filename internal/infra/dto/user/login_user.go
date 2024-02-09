package userdto

import (
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type LoginUserInput struct {
	userentity.UserCommonAttributes
}

func (u *LoginUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if !utils.IsValidPassword(u.Password) {
		return ErrPasswordInvalid
	}

	return nil
}

func (u *LoginUserInput) ToModel() (*userentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	userCommonAttributes := userentity.UserCommonAttributes{
		Email:    u.Email,
		Password: u.Password,
	}

	return userentity.NewUser(userCommonAttributes), nil
}
