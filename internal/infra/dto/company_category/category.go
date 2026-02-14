package companycategorydto

import (
	"github.com/google/uuid"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
)

type CreateCategoryDTO struct {
	Name      string `json:"name" validate:"required"`
	ImagePath string `json:"image_path"`
}

func (dto *CreateCategoryDTO) ToEntity() companycategoryentity.CompanyCategoryCommonAttributes {
	return companycategoryentity.CompanyCategoryCommonAttributes{
		Name:      dto.Name,
		ImagePath: dto.ImagePath,
	}
}

type UpdateCategoryDTO struct {
	Name      string `json:"name"`
	ImagePath string `json:"image_path"`
}

func (dto *UpdateCategoryDTO) ToEntity() companycategoryentity.CompanyCategoryCommonAttributes {
	return companycategoryentity.CompanyCategoryCommonAttributes{
		Name:      dto.Name,
		ImagePath: dto.ImagePath,
	}
}

type CompanyCategoryDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ImagePath string    `json:"image_path"`
}

func (dto *CompanyCategoryDTO) FromDomain(category *companycategoryentity.CompanyCategory) {
	dto.ID = category.ID
	dto.Name = category.Name
	dto.ImagePath = category.ImagePath
}
