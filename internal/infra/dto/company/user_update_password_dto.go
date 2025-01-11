package companydto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrMustBeDifferentPassword = errors.New("must be different password")
)

type UserUpdatePasswordDTO struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (r *UserUpdatePasswordDTO) validate() error {
	if !utils.IsEmailValid(r.Email) {
		return ErrEmailInvalid
	}

	if r.CurrentPassword == r.NewPassword {
		return ErrMustBeDifferentPassword
	}

	if err := utils.ValidatePassword(r.NewPassword); err != nil {
		return err
	}

	return nil
}

func (r *UserUpdatePasswordDTO) ToDomain() (*companyentity.User, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Email: r.Email,
	}

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person:   *person,
		Password: r.CurrentPassword,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
