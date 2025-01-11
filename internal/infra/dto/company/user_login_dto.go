package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLoginDTO) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if err := utils.ValidatePassword(u.Password); err != nil && u.Password != "12345" {
		return err
	}

	return nil
}

func (u *UserLoginDTO) ToDomain() (*companyentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Email: u.Email,
	}

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person:   *person,
		Password: u.Password,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
