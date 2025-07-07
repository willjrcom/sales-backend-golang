package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type UserUpdateForgetPasswordDTO struct {
	Token    string `json:"token"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserUpdateForgetPasswordDTO) validate() error {
	if !utils.IsEmailValid(r.Email) {
		return ErrEmailInvalid
	}

	if err := utils.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

func (r *UserUpdateForgetPasswordDTO) ToDomain() (*companyentity.User, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Email: r.Email,
	}

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person:   *person,
		Password: r.Password,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
