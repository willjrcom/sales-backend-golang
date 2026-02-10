package companydto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

var (
	ErrEmailInvalid             = errors.New("email is invalid")
	ErrMustHaveAtLeastOneSchema = errors.New("must have at least one schema")
)

type UserCreateDTO struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	GeneratePassword bool   `json:"generate_password"`
	Name             string `json:"name"`
	Cpf              string `json:"cpf,omitempty"`
}

func (u *UserCreateDTO) validate() error {
	if !utils.IsEmailValid(u.Email) {
		return ErrEmailInvalid
	}

	if err := utils.ValidatePassword(u.Password); err != nil && !u.GeneratePassword {
		return err
	}

	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Cpf == "" {
		return errors.New("cpf is required")
	}
	return nil
}

func (u *UserCreateDTO) ToDomain() (*companyentity.User, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	if u.GeneratePassword {
		u.Password = utils.GeneratePassword(10, true, true, true)
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Name:  u.Name,
		Email: u.Email,
		Cpf:   u.Cpf,
	}

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person:   *person,
		Password: u.Password,
	}
	user := companyentity.NewUser(userCommonAttributes)

	return user, nil
}
