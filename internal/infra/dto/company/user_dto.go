package companydto

import (
	"time"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type UserDTO struct {
	ID        uuid.UUID              `json:"id"`
	Email     string                 `json:"email"`
	Name      string                 `json:"name"`
	Cpf       string                 `json:"cpf,omitempty"`
	Birthday  *time.Time             `json:"birthday,omitempty"`
	Contact   *contactdto.ContactDTO `json:"contact,omitempty"`
	Address   *addressdto.AddressDTO `json:"address,omitempty"`
	Companies []CompanyDTO           `json:"companies,omitempty"`
}

func (u *UserDTO) FromDomain(user *companyentity.User) {
	if user == nil {
		return
	}
	*u = UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Cpf:       user.Cpf,
		Birthday:  user.Birthday,
		Contact:   &contactdto.ContactDTO{},
		Address:   &addressdto.AddressDTO{},
		Companies: []CompanyDTO{},
	}

	u.Contact.FromDomain(user.Contact)
	u.Address.FromDomain(user.Address)

	for _, company := range user.Companies {
		companyDTO := CompanyDTO{}
		companyDTO.FromDomain(&company)
		u.Companies = append(u.Companies, companyDTO)
	}

	if user.Contact == nil {
		u.Contact = nil
	}
	if user.Address == nil {
		u.Address = nil
	}
	if len(user.Companies) == 0 {
		u.Companies = nil
	}
}
