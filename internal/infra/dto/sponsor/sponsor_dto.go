package sponsordto

import (
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CreateSponsorDTO struct {
	Name        string                 `json:"name"`
	CNPJ        string                 `json:"cnpj"`
	Email       string                 `json:"email"`
	Contact     string                 `json:"contact"`
	Address     *addressdto.AddressDTO `json:"address"`
	CategoryIDs []string               `json:"category_ids,omitempty"`
}

func (dto *CreateSponsorDTO) ToDomain() (*sponsorentity.Sponsor, error) {
	cats := make([]sponsorentity.CompanyCategory, len(dto.CategoryIDs))
	for i, id := range dto.CategoryIDs {
		cats[i] = sponsorentity.CompanyCategory{ID: id}
	}

	attr := sponsorentity.SponsorCommonAttributes{
		Name:                   dto.Name,
		CNPJ:                   dto.CNPJ,
		Email:                  dto.Email,
		Contact:                dto.Contact,
		CompanyCategorySponsor: cats,
	}

	if dto.Address != nil {
		address, err := dto.Address.ToDomain()
		if err != nil {
			return nil, err
		}
		attr.Address = address
	}

	return sponsorentity.NewSponsor(attr), nil
}

type UpdateSponsorDTO struct {
	Name        *string                `json:"name,omitempty"`
	CNPJ        *string                `json:"cnpj,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Contact     *string                `json:"contact,omitempty"`
	Address     *addressdto.AddressDTO `json:"address,omitempty"`
	CategoryIDs []string               `json:"category_ids,omitempty"`
}

func (dto *UpdateSponsorDTO) UpdateDomain(s *sponsorentity.Sponsor) error {
	if dto.Name != nil {
		s.Name = *dto.Name
	}
	if dto.CNPJ != nil {
		s.CNPJ = *dto.CNPJ
	}
	if dto.Email != nil {
		s.Email = *dto.Email
	}
	if dto.Contact != nil {
		s.Contact = *dto.Contact
	}
	if dto.CategoryIDs != nil {
		cats := make([]sponsorentity.CompanyCategory, len(dto.CategoryIDs))
		for i, id := range dto.CategoryIDs {
			cats[i] = sponsorentity.CompanyCategory{ID: id}
		}
		s.CompanyCategorySponsor = cats
	}
	if dto.Address != nil {
		if s.Address == nil {
			address, err := dto.Address.ToDomain()
			if err != nil {
				return err
			}
			s.Address = address
		} else {
			if err := dto.Address.UpdateDomain(s.Address); err != nil {
				return err
			}
		}
	}
	return nil
}
