package companyrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyRepositoryLocal struct {
	companies          map[uuid.UUID]*model.Company
	publicCompanyUsers map[uuid.UUID]struct{}
	mu                 sync.RWMutex
}

func NewCompanyRepositoryLocal() model.CompanyRepository {
	return &CompanyRepositoryLocal{
		companies:          make(map[uuid.UUID]*model.Company),
		publicCompanyUsers: make(map[uuid.UUID]struct{}),
	}
}

func (r *CompanyRepositoryLocal) NewCompany(ctx context.Context, company *model.Company) error {
	if company == nil || company.Entity.ID == uuid.Nil {
		return errors.New("invalid company")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.companies[company.Entity.ID]; exists {
		return errors.New("company already exists")
	}
	r.companies[company.Entity.ID] = company
	return nil
}

func (r *CompanyRepositoryLocal) UpdateCompany(ctx context.Context, company *model.Company) error {
	if company == nil || company.Entity.ID == uuid.Nil {
		return errors.New("invalid company")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.companies[company.Entity.ID]; !exists {
		return errors.New("company not found")
	}
	r.companies[company.Entity.ID] = company
	return nil
}

func (r *CompanyRepositoryLocal) GetCompany(ctx context.Context) (*model.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.companies {
		return c, nil // return first company found
	}
	return nil, errors.New("no company found")
}

func (r *CompanyRepositoryLocal) ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.publicCompanyUsers[userID]
	return exists, nil
}

func (r *CompanyRepositoryLocal) AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("invalid userID")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.publicCompanyUsers[userID] = struct{}{}
	return nil
}

func (r *CompanyRepositoryLocal) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("invalid userID")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.publicCompanyUsers, userID)
	return nil
}
