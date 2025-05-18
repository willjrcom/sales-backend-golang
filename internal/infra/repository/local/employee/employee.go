package employeerepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errEmployeeExists   = errors.New("employee already exists")
	errEmployeeNotFound = errors.New("employee not found")
)

// EmployeeRepositoryLocal is an in-memory implementation of EmployeeRepository
type EmployeeRepositoryLocal struct {
	mu        sync.RWMutex
	employees map[uuid.UUID]*model.Employee
	payments  map[string][]model.PaymentEmployee
}

func NewEmployeeRepositoryLocal() model.EmployeeRepository {
	return &EmployeeRepositoryLocal{
		employees: make(map[uuid.UUID]*model.Employee),
		payments:  make(map[string][]model.PaymentEmployee),
	}
}

func (r *EmployeeRepositoryLocal) CreateEmployee(_ context.Context, p *model.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.employees[p.ID]; ok {
		return errEmployeeExists
	}
	r.employees[p.ID] = p
	return nil
}

func (r *EmployeeRepositoryLocal) UpdateEmployee(_ context.Context, p *model.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.employees[p.ID] = p
	return nil
}

func (r *EmployeeRepositoryLocal) DeleteEmployee(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	uid := uuid.MustParse(id)
	if _, ok := r.employees[uid]; !ok {
		return errEmployeeNotFound
	}
	delete(r.employees, uid)
	return nil
}

func (r *EmployeeRepositoryLocal) GetEmployeeById(_ context.Context, id string) (*model.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	uid := uuid.MustParse(id)
	if e, ok := r.employees[uid]; ok {
		return e, nil
	}
	return nil, errEmployeeNotFound
}

func (r *EmployeeRepositoryLocal) GetEmployeeByUserID(_ context.Context, userID string) (*model.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	uid := uuid.MustParse(userID)
	for _, e := range r.employees {
		if e.UserID == uid {
			return e, nil
		}
	}
	return nil, errEmployeeNotFound
}

func (r *EmployeeRepositoryLocal) GetAllEmployees(_ context.Context) ([]model.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Employee, 0, len(r.employees))
	for _, e := range r.employees {
		list = append(list, *e)
	}
	return list, nil
}

// AddPaymentEmployee records a payment for an employee in memory
func (r *EmployeeRepositoryLocal) AddPaymentEmployee(_ context.Context, p *model.PaymentEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := p.EmployeeID.String()
	// ensure employee exists
	if _, ok := r.employees[p.EmployeeID]; !ok {
		return errEmployeeNotFound
	}
	r.payments[key] = append(r.payments[key], *p)
	return nil
}
