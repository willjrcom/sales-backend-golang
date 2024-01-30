package userdto

import (
	"errors"

	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrEmailInvalid             = errors.New("email is invalid")
	ErrPasswordInvalid          = errors.New("password is invalid")
	ErrEmployeeIDRequired       = errors.New("employee id is required")
	ErrMustHaveAtLeastOneSchema = errors.New("must have at least one schema")
)

type CreateUserInput struct {
	userentity.UserCommonAttributes
}

func (u *CreateUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if !utils.IsValidPassword(u.Password) {
		return ErrPasswordInvalid
	}

	if u.EmployeeID == nil {
		return ErrEmployeeIDRequired
	}

	if len(u.Schemas) == 0 {
		return ErrMustHaveAtLeastOneSchema
	}

	return nil
}

func (u *CreateUserInput) ToModel() (*userentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	return userentity.NewUser(u.UserCommonAttributes), nil
}
