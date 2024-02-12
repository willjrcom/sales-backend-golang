package userdto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrMustBeDifferentPassword = errors.New("must be different password")
)

type UpdatePasswordInput struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (r *UpdatePasswordInput) validate() error {
	if !utils.IsEmailValid(r.Email) {
		return ErrEmailInvalid
	}

	if r.OldPassword == r.NewPassword {
		return ErrMustBeDifferentPassword
	}

	if !utils.IsValidPassword(r.NewPassword) {
		return ErrPasswordInvalid
	}

	return nil
}

func (r *UpdatePasswordInput) ToModel() (*companyentity.User, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	userCommonAttributes := companyentity.UserCommonAttributes{
		Email:    r.Email,
		Password: r.OldPassword,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
