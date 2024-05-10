package companycategoryentity

import "context"

type CompanyCategoryRepository interface {
	CreateCompanyCategory(ctx context.Context, CompanyCategory *CompanyCategory) (err error)
	UpdateCompanyCategory(ctx context.Context, CompanyCategory *CompanyCategory) (err error)
	DeleteCompanyCategory(ctx context.Context, id string) (err error)
	GetCompanyCategoryByID(ctx context.Context, id string) (CompanyCategory *CompanyCategory, err error)
	GetAllCompanyCategories(ctx context.Context) ([]CompanyCategory, error)
}
