package userdto

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type ForgetUserPassword struct {
	Email string `json:"email"`
}

func (r *ForgetUserPassword) validate() error {
	if !utils.IsEmailValid(r.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (r *ForgetUserPassword) ToModel() (*string, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	return &r.Email, nil
}
