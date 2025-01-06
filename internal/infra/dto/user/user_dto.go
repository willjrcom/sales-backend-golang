package userdto

import (
	"time"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type UserDTO struct {
	ID       uuid.UUID              `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Cpf      string                 `json:"cpf,omitempty"`
	Birthday *time.Time             `json:"birthday,omitempty"`
	Contact  *contactdto.ContactDTO `json:"contact,omitempty"`
	Address  *addressdto.AddressDTO `json:"address,omitempty"`
}

func (u *UserDTO) FromModel(user *companyentity.User) {
	*u = UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Cpf:      user.Cpf,
		Birthday: user.Birthday,
	}

	u.Contact.FromDomain(user.Contact)
	u.Address.FromModel(user.Address)
}
