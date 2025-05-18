package companyrepositorylocal

import (
	"context"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyRepositoryLocal struct {}

func NewCompanyRepositoryLocal() model.CompanyRepository {
	return &CompanyRepositoryLocal{}
}

func (r *CompanyRepositoryLocal) NewCompany(ctx context.Context, company *model.Company) error {
	return nil
}

func (r *CompanyRepositoryLocal) UpdateCompany(ctx context.Context, company *model.Company) error {
	return nil
}

func (r *CompanyRepositoryLocal) GetCompany(ctx context.Context) (*model.Company, error) {
	return nil, nil
}

func (r *CompanyRepositoryLocal) ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error) {
	return false, nil
}

func (r *CompanyRepositoryLocal) AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (r *CompanyRepositoryLocal) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	return nil
}
