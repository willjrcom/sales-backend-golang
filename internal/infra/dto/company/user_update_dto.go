package companydto

import (
	"errors"
	"strings"
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type UserUpdateDTO struct {
	Name      *string                      `json:"name"`
	Email     *string                      `json:"email"`
	ImagePath *string                      `json:"image_path"`
	Cpf       *string                      `json:"cpf"`
	Birthday  *time.Time                   `json:"birthday"`
	Contact   *contactdto.ContactUpdateDTO `json:"contact"`
	Address   *addressdto.AddressUpdateDTO `json:"address"`
}

func (r *UserUpdateDTO) validate() error {
	if r.Email != nil && !strings.Contains(*r.Email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (r *UserUpdateDTO) UpdateDomain(user *companyentity.User) error {
	if err := r.validate(); err != nil {
		return err
	}

	if r.Name != nil {
		user.Name = *r.Name
	}
	if r.ImagePath != nil {
		user.ImagePath = *r.ImagePath
	}
	if r.Birthday != nil {
		user.Birthday = r.Birthday
	}
	if r.Contact != nil {
		if user.Contact == nil {
			user.Contact = &personentity.Contact{
				Entity: entity.NewEntity(),
				ContactCommonAttributes: personentity.ContactCommonAttributes{
					Number: r.Contact.Number,
				},
				ObjectID: user.ID,
			}
		}
		r.Contact.UpdateDomain(user.Contact, personentity.ContactTypeEmployee)
	} else {
		user.Contact = nil
	}

	if r.Address != nil {
		if user.Address == nil {
			user.Address = &addressentity.Address{
				Entity:                  entity.NewEntity(),
				AddressCommonAttributes: addressentity.AddressCommonAttributes{},
				ObjectID:                user.ID,
			}
		}
		r.Address.UpdateDomain(user.Address)
	} else {
		user.Address = nil
	}

	return nil
}
