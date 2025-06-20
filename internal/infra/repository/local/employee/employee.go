package employeerepositorylocal

import (
	"context"
	"errors"
	"sort"
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
}

func NewEmployeeRepositoryLocal() model.EmployeeRepository {
	return &EmployeeRepositoryLocal{
		employees: make(map[uuid.UUID]*model.Employee),
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

func (r *EmployeeRepositoryLocal) GetEmployeeDeletedByUserID(_ context.Context, userID string) (*model.Employee, error) {
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

// GetAllEmployees retrieves a paginated list of employees and the total count.
func (r *EmployeeRepositoryLocal) GetAllEmployees(_ context.Context, page, perPage int) ([]model.Employee, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.employees)
	if total == 0 {
		return []model.Employee{}, 0, nil
	}
	ids := make([]uuid.UUID, 0, total)
	for id := range r.employees {
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
		return []model.Employee{}, total, nil
	}
	end := offset + perPage
	if end > total {
		end = total
	}
	segment := ids[offset:end]
	result := make([]model.Employee, 0, len(segment))
	for _, id := range segment {
		result = append(result, *r.employees[id])
	}
	return result, total, nil
}

// GetAllEmployees retrieves a paginated list of employees and the total count.
func (r *EmployeeRepositoryLocal) GetAllEmployeeDeleted(_ context.Context, page, perPage int) ([]model.Employee, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.employees)
	if total == 0 {
		return []model.Employee{}, 0, nil
	}
	ids := make([]uuid.UUID, 0, total)
	for id := range r.employees {
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
		return []model.Employee{}, total, nil
	}
	end := offset + perPage
	if end > total {
		end = total
	}
	segment := ids[offset:end]
	result := make([]model.Employee, 0, len(segment))
	for _, id := range segment {
		result = append(result, *r.employees[id])
	}
	return result, total, nil
}

// AddPaymentEmployee records a payment for an employee in memory
func (r *EmployeeRepositoryLocal) AddPaymentEmployee(_ context.Context, p *model.PaymentEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ensure employee exists
	employee, ok := r.employees[p.EmployeeID]
	if !ok {
		return errEmployeeNotFound
	}

	employee.Payments = append(employee.Payments, *p)
	r.employees[p.EmployeeID] = employee
	return nil
}

func (r *EmployeeRepositoryLocal) GetSalaryHistory(_ context.Context, employeeID uuid.UUID) ([]model.EmployeeSalaryHistory, error) {
	// Retorna slice vazio para ambiente local (mock)
	return []model.EmployeeSalaryHistory{}, nil
}

func (r *EmployeeRepositoryLocal) GetPayments(_ context.Context, employeeID uuid.UUID) ([]model.PaymentEmployee, error) {
	// Retorna slice vazio para ambiente local (mock)
	return []model.PaymentEmployee{}, nil
}
