package companycategoryusecases

import (
	"context"

	"github.com/google/uuid"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	repo model.CompanyCategoryRepository
}

func NewService(repo model.CompanyCategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, attrs companycategoryentity.CompanyCategoryCommonAttributes) (*companycategoryentity.CompanyCategory, error) {
	category := companycategoryentity.NewCategory(attrs)

	categoryModel := &model.CompanyCategory{}
	categoryModel.FromDomain(category)

	if err := s.repo.Create(ctx, categoryModel); err != nil {
		return nil, err
	}

	return categoryModel.ToDomain(), nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uuid.UUID, attrs companycategoryentity.CompanyCategoryCommonAttributes) (*companycategoryentity.CompanyCategory, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = attrs.Name
	category.ImagePath = attrs.ImagePath

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category.ToDomain(), nil
}

func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCategory(ctx context.Context, id uuid.UUID) (*companycategoryentity.CompanyCategory, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return category.ToDomain(), nil
}

func (s *Service) GetAllCompanyCategories(ctx context.Context) ([]companycategoryentity.CompanyCategory, error) {
	categories, err := s.repo.GetAllCompanyCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]companycategoryentity.CompanyCategory, len(categories))
	for i, cat := range categories {
		result[i] = *cat.ToDomain()
	}

	return result, nil
}
