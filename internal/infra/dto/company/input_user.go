package companydto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

var (
	ErrMustBeEmail    = errors.New("email is required")
	ErrMustBePassword = errors.New("password is required")
)

type UserInput struct {
	companyentity.UserCommonAttributes
}

func (c *UserInput) validate() error {
	if c.Email == "" {
		return ErrMustBeEmail
	}

	return nil
}

func (c *UserInput) ToModel() (*companyentity.User, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return companyentity.NewUser(c.UserCommonAttributes), nil
}
