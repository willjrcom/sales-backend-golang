package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

// UserBasicDTO exposes minimal user identification info.
type UserBasicDTO struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Cpf       string                 `json:"cpf"`
	Contact   *contactdto.ContactDTO `json:"contact,omitempty"`
	Companies []CompanyDTO           `json:"companies"`
}

func (u *UserBasicDTO) FromDomain(user *companyentity.User) {
	if user == nil {
		return
	}

	*u = UserBasicDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Cpf:       user.Cpf,
		Companies: []CompanyDTO{},
	}

	if user.Contact != nil {
		u.Contact = &contactdto.ContactDTO{}
		u.Contact.FromDomain(user.Contact)
	}

	for _, company := range user.Companies {
		companyDTO := CompanyDTO{}
		companyDTO.FromDomain(&company)
		u.Companies = append(u.Companies, companyDTO)
	}
}
