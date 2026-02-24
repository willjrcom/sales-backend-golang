package advertisingdto

import (
	"github.com/google/uuid"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
	companycategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company_category"
	sponsordto "github.com/willjrcom/sales-backend-go/internal/infra/dto/sponsor"
)

type CreateAdvertisingDTO struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Link           string   `json:"link"`
	Contact        string   `json:"contact"`
	CoverImagePath string   `json:"cover_image_path"`
	Images         []string `json:"images"`
	SponsorID      string   `json:"sponsor_id"`
	CategoryIDs    []string `json:"category_ids,omitempty"`
}

func (dto *CreateAdvertisingDTO) ToDomain() (*advertisingentity.Advertising, error) {
	sponsorID, err := uuid.Parse(dto.SponsorID)
	if err != nil {
		return nil, err
	}

	cats := make([]advertisingentity.CompanyCategory, len(dto.CategoryIDs))
	for i, id := range dto.CategoryIDs {
		cats[i] = advertisingentity.CompanyCategory{ID: id}
	}

	attr := advertisingentity.AdvertisingCommonAttributes{
		Title:                      dto.Title,
		Description:                dto.Description,
		Link:                       dto.Link,
		Contact:                    dto.Contact,
		CoverImagePath:             dto.CoverImagePath,
		Images:                     dto.Images,
		SponsorID:                  sponsorID,
		CompanyCategoryAdvertising: cats,
	}

	return advertisingentity.NewAdvertising(attr), nil
}

type UpdateAdvertisingDTO struct {
	Title          *string  `json:"title,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Link           *string  `json:"link,omitempty"`
	Contact        *string  `json:"contact,omitempty"`
	CoverImagePath *string  `json:"cover_image_path,omitempty"`
	Images         []string `json:"images,omitempty"`
	SponsorID      *string  `json:"sponsor_id,omitempty"`
	CategoryIDs    []string `json:"category_ids,omitempty"`
}

func (dto *UpdateAdvertisingDTO) UpdateDomain(a *advertisingentity.Advertising) error {
	if dto.Title != nil {
		a.Title = *dto.Title
	}
	if dto.Description != nil {
		a.Description = *dto.Description
	}
	if dto.Link != nil {
		a.Link = *dto.Link
	}
	if dto.Contact != nil {
		a.Contact = *dto.Contact
	}
	if dto.CoverImagePath != nil {
		a.CoverImagePath = *dto.CoverImagePath
	}
	if dto.Images != nil {
		a.Images = dto.Images
	}
	if dto.SponsorID != nil {
		sponsorID, err := uuid.Parse(*dto.SponsorID)
		if err != nil {
			return err
		}
		a.SponsorID = sponsorID
	}
	if dto.CategoryIDs != nil {
		cats := make([]advertisingentity.CompanyCategory, len(dto.CategoryIDs))
		for i, id := range dto.CategoryIDs {
			cats[i] = advertisingentity.CompanyCategory{ID: id}
		}
		a.CompanyCategoryAdvertising = cats
	}
	return nil
}

type AdvertisingDTO struct {
	ID             uuid.UUID                               `json:"id"`
	Title          string                                  `json:"title"`
	Description    string                                  `json:"description"`
	Link           string                                  `json:"link"`
	Contact        string                                  `json:"contact"`
	CoverImagePath string                                  `json:"cover_image_path"`
	Images         []string                                `json:"images"`
	SponsorID      uuid.UUID                               `json:"sponsor_id"`
	Sponsor        *sponsordto.SponsorDTO                  `json:"sponsor"`
	Categories     []companycategorydto.CompanyCategoryDTO `json:"categories"`
}

func (dto *AdvertisingDTO) FromDomain(a *advertisingentity.Advertising) {
	dto.ID = a.ID
	dto.Title = a.Title
	dto.Description = a.Description
	dto.Link = a.Link
	dto.Contact = a.Contact
	dto.CoverImagePath = a.CoverImagePath
	dto.Images = a.Images
	dto.SponsorID = a.SponsorID

	if a.Sponsor != nil {
		dto.Sponsor = &sponsordto.SponsorDTO{}
		dto.Sponsor.FromDomain(a.Sponsor)
	}

	dto.Categories = make([]companycategorydto.CompanyCategoryDTO, len(a.CompanyCategoryAdvertising))
	for i, cat := range a.CompanyCategoryAdvertising {
		id, _ := uuid.Parse(cat.ID)
		dto.Categories[i] = companycategorydto.CompanyCategoryDTO{
			ID: id,
		}
	}
}
