package model

import (
	"context"

	"github.com/google/uuid"
)

type CompanyCategoryRepository interface {
	Create(ctx context.Context, category *CompanyCategory) error
	Update(ctx context.Context, category *CompanyCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*CompanyCategory, error)
	GetAllCompanyCategories(ctx context.Context) ([]CompanyCategory, error)
}
