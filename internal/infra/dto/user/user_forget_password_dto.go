package userdto

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/service/utils"
)

type UserForgetPasswordDTO struct {
	Email string `json:"email"`
}

func (r *UserForgetPasswordDTO) validate() error {
	if !utils.IsEmailValid(r.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (r *UserForgetPasswordDTO) ToModel() (*string, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	return &r.Email, nil
}
