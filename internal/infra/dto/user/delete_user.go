package userdto

import (
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type DeleteUserInput struct {
	userentity.UserCommonAttributes
}

func (u *DeleteUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if !utils.IsValidPassword(u.Password) {
		return ErrPasswordInvalid
	}

	return nil
}

func (u *DeleteUserInput) ToModel() (*userentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	userCommonAttributes := userentity.UserCommonAttributes{
		Email:    u.Email,
		Password: u.Password,
	}

	return userentity.NewUser(userCommonAttributes), nil
}
