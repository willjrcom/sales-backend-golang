package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type UserDeleteDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserDeleteDTO) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return err
	}

	return nil
}

func (u *UserDeleteDTO) ToDomain() (*companyentity.User, error) {
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
