package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type UserBasicCreateDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserBasicCreateDTO) validate() error {
	if c.Email == "" {
		return ErrMustBeEmail
	}
	if c.Password == "" {
		return ErrMustBePassword
	}
	return nil
}

func (c *UserBasicCreateDTO) ToDomain() (*companyentity.User, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	personCommonAttributes := &personentity.PersonCommonAttributes{
		Email: c.Email,
	}

	person := personentity.NewPerson(personCommonAttributes)

	userCommonAttributes := &companyentity.UserCommonAttributes{
		Person: *person,
	}

	return companyentity.NewUser(userCommonAttributes), nil
}
