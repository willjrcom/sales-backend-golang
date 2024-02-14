package userdto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrEmailInvalid             = errors.New("email is invalid")
	ErrPasswordInvalid          = errors.New("password is invalid")
	ErrMustHaveAtLeastOneSchema = errors.New("must have at least one schema")
)

type CreateUserInput struct {
	companyentity.UserCommonAttributes
	GeneratePassword bool `json:"generate_password"`
}

func (u *CreateUserInput) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if !utils.IsValidPassword(u.Password) && !u.GeneratePassword {
		return ErrPasswordInvalid
	}

	return nil
}

func (u *CreateUserInput) ToModel() (*companyentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	if u.GeneratePassword {
		u.Password = "12345"
	}

	return companyentity.NewUser(u.UserCommonAttributes), nil
}
