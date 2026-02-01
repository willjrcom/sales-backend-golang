package companyrepositorylocal

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyRepositoryLocal struct {
	companies          map[uuid.UUID]*model.Company
	publicCompanyUsers map[uuid.UUID]struct{}
	payments           map[string]*model.CompanyPayment
	mu                 sync.RWMutex
}

func NewCompanyRepositoryLocal() model.CompanyRepository {
	return &CompanyRepositoryLocal{
		companies:          make(map[uuid.UUID]*model.Company),
		publicCompanyUsers: make(map[uuid.UUID]struct{}),
		payments:           make(map[string]*model.CompanyPayment),
	}
}

func (r *CompanyRepositoryLocal) NewCompany(ctx context.Context, company *model.Company) error {
	if company == nil || company.Entity.ID == uuid.Nil {
		return errors.New("invalid company")
	}

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

	if _, exists := r.companies[company.Entity.ID]; !exists {
		return errors.New("company not found")
	}
	r.companies[company.Entity.ID] = company
	return nil
}

func (r *CompanyRepositoryLocal) GetCompany(ctx context.Context, withouRelations ...bool) (*model.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.companies {
		return c, nil // return first company found
	}
	return nil, errors.New("no company found")
}

func (r *CompanyRepositoryLocal) ListPublicCompanies(ctx context.Context) ([]model.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Company, 0, len(r.companies))
	for _, company := range r.companies {
		result = append(result, *company)
	}
	return result, nil
}

func (r *CompanyRepositoryLocal) ListCompaniesForBilling(ctx context.Context) ([]model.Company, error) {
	return r.ListPublicCompanies(ctx)
}

func (r *CompanyRepositoryLocal) ListCompaniesByPaymentDueDay(ctx context.Context, day int) ([]model.Company, error) {
	// Stub implementation
	// In a real local scenario, we would filter r.companies
	// For now, returning empty to satisfy interface
	return []model.Company{}, nil
}

func (r *CompanyRepositoryLocal) UpdateCompanySubscription(ctx context.Context, companyID uuid.UUID, schema string, expiresAt *time.Time, planType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	company, ok := r.companies[companyID]
	if !ok {
		return fmt.Errorf("company %s not found", companyID)
	}
	company.SubscriptionExpiresAt = expiresAt
	return nil
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

	r.publicCompanyUsers[userID] = struct{}{}
	return nil
}

func (r *CompanyRepositoryLocal) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("invalid userID")
	}

	delete(r.publicCompanyUsers, userID)
	return nil
}

// GetCompanyUsers retrieves a paginated list of users in the public company and the total count.
func (r *CompanyRepositoryLocal) GetCompanyUsers(ctx context.Context, page, perPage int) ([]model.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.publicCompanyUsers)
	if total == 0 {
		return []model.User{}, 0, nil
	}
	ids := make([]uuid.UUID, 0, total)
	for id := range r.publicCompanyUsers {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].String() < ids[j].String()
	})
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = total
	}
	offset := (page - 1) * perPage
	if offset >= total {
		return []model.User{}, total, nil
	}
	end := offset + perPage
	if end > total {
		end = total
	}
	segment := ids[offset:end]
	result := make([]model.User, 0, len(segment))
	for _, id := range segment {
		u := model.User{}
		u.Entity.ID = id
		result = append(result, u)
	}
	return result, total, nil
}

func (r *CompanyRepositoryLocal) CreateCompanyPayment(ctx context.Context, payment *model.CompanyPayment) error {
	if payment == nil {
		return errors.New("payment is nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := payment.Provider + ":" + payment.ProviderPaymentID
	if _, exists := r.payments[key]; exists {
		return fmt.Errorf("payment %s already exists", payment.ProviderPaymentID)
	}
	r.payments[key] = payment
	return nil
}

func (r *CompanyRepositoryLocal) UpdateCompanyPayment(ctx context.Context, payment *model.CompanyPayment) error {
	if payment == nil || payment.ID == uuid.Nil {
		return errors.New("invalid payment")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update index by provider if changed?
	// For simplicity, we just store it. But if provider ID changed, we should re-index.
	// Since we use key provider:providerID, we should handle that.
	key := payment.Provider + ":" + payment.ProviderPaymentID
	r.payments[key] = payment
	return nil
}

func (r *CompanyRepositoryLocal) GetCompanyPaymentByID(ctx context.Context, id uuid.UUID) (*model.CompanyPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Since payments are indexed by provider+id, we just scan.
	// Not efficient but acceptable for local stub.
	for _, p := range r.payments {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("payment not found")
}

func (r *CompanyRepositoryLocal) ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]model.CompanyPayment, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matching := make([]*model.CompanyPayment, 0)
	for _, payment := range r.payments {
		if payment.CompanyID == companyID {
			matching = append(matching, payment)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		return matching[i].PaidAt.After(*matching[j].PaidAt)
	})

	total := len(matching)
	if total == 0 {
		return []model.CompanyPayment{}, 0, nil
	}

	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = total
	}

	offset := (page - 1) * perPage
	if offset >= total {
		return []model.CompanyPayment{}, total, nil
	}

	end := offset + perPage
	if end > total {
		end = total
	}

	result := make([]model.CompanyPayment, 0, end-offset)
	for _, payment := range matching[offset:end] {
		result = append(result, *payment)
	}

	return result, total, nil
}

func (r *CompanyRepositoryLocal) UpdateBlockStatus(ctx context.Context, companyID uuid.UUID, isBlocked bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	company, ok := r.companies[companyID]
	if !ok {
		return fmt.Errorf("company %s not found", companyID)
	}
	company.IsBlocked = isBlocked
	return nil
}

func (r *CompanyRepositoryLocal) ListOverduePayments(ctx context.Context, cutoffDate time.Time) ([]model.CompanyPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var overdue []model.CompanyPayment
	for _, p := range r.payments {
		if p.Status == "pending" && p.IsMandatory && p.ExpiresAt != nil && p.ExpiresAt.Before(cutoffDate) {
			overdue = append(overdue, *p)
		}
	}
	return overdue, nil
}

func (r *CompanyRepositoryLocal) ListPendingMandatoryPayments(ctx context.Context, companyID uuid.UUID) ([]model.CompanyPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var pending []model.CompanyPayment
	for _, p := range r.payments {
		if p.CompanyID == companyID && p.Status == "pending" && p.IsMandatory {
			pending = append(pending, *p)
		}
	}
	return pending, nil
}

func (r *CompanyRepositoryLocal) CreateSubscription(ctx context.Context, subscription *model.CompanySubscription) error {
	return nil // Stub
}

func (r *CompanyRepositoryLocal) UpdateSubscription(ctx context.Context, subscription *model.CompanySubscription) error {
	return nil // Stub
}

func (r *CompanyRepositoryLocal) MarkActiveSubscriptionAsCanceled(ctx context.Context, companyID uuid.UUID) error {
	return nil // Stub
}

func (r *CompanyRepositoryLocal) GetActiveSubscription(ctx context.Context, companyID uuid.UUID) (*model.CompanySubscription, error) {
	return nil, nil // Stub
}

func (r *CompanyRepositoryLocal) GetUpcomingSubscription(ctx context.Context, companyID uuid.UUID) (*model.CompanySubscription, error) {
	return nil, nil // Stub
}

func (r *CompanyRepositoryLocal) UpdateCompanyPlans(ctx context.Context) error {
	return nil // Stub
}
