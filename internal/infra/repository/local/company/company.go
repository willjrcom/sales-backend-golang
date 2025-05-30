package companyrepositorylocal

import (
   "context"
   "errors"
   "sort"
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

// GetCompanyUsers retrieves a paginated list of users in the public company and the total count.
func (r *CompanyRepositoryLocal) GetCompanyUsers(ctx context.Context, offset, limit int) ([]model.User, int, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   total := len(r.publicCompanyUsers)
   if total == 0 {
       return []model.User{}, 0, nil
   }
   // collect and sort user IDs for consistent ordering
   ids := make([]uuid.UUID, 0, total)
   for id := range r.publicCompanyUsers {
       ids = append(ids, id)
   }
   sort.Slice(ids, func(i, j int) bool {
       return ids[i].String() < ids[j].String()
   })
   // normalize pagination parameters
   if offset < 0 {
       offset = 0
   }
   if limit <= 0 {
       limit = total
   }
   if offset > total {
       offset = total
   }
   end := offset + limit
   if end > total {
       end = total
   }
   // build result slice
   users := make([]model.User, 0, end-offset)
   for _, id := range ids[offset:end] {
       u := model.User{}
       u.Entity.ID = id
       users = append(users, u)
   }
   return users, total, nil
}
