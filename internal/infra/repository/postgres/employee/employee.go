package employeerepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type EmployeeRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewEmployeeRepositoryBun(db *bun.DB) model.EmployeeRepository {
	return &EmployeeRepositoryBun{db: db}
}

func (r *EmployeeRepositoryBun) CreateEmployee(ctx context.Context, c *model.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	// Create employee
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepositoryBun) UpdateEmployee(ctx context.Context, p *model.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(p).Where("employee.id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	// Delete employee
	if _, err := tx.NewDelete().Model(&model.Employee{}).Where("employee.id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*model.Employee, error) {
	employee := &model.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(employee).Where("employee.id = ?", id).Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *EmployeeRepositoryBun) GetEmployeeByUserID(ctx context.Context, userID string) (*model.Employee, error) {
	employee := &model.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(employee).Where("employee.user_id = ?", userID).Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *EmployeeRepositoryBun) GetEmployeeDeletedByUserID(ctx context.Context, userID string) (*model.Employee, error) {
	employee := &model.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(employee).Where("employee.user_id = ?", userID).WhereAllWithDeleted().Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

// GetAllEmployees retrieves a paginated list of employees and the total count.
func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context, page, perPage int) ([]model.Employee, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}
	// count total records
	totalCount, err := tx.NewSelect().Model((*model.Employee)(nil)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	// fetch paginated records
	employees := []model.Employee{}
	err = tx.NewSelect().
		Model(&employees).
		Relation("User").
		Relation("User.Address").
		Relation("User.Contact").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return employees, int(totalCount), nil
}

func (r *EmployeeRepositoryBun) AddPaymentEmployee(ctx context.Context, p *model.PaymentEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetAllEmployeeDeleted retrieves a paginated list of soft-deleted employees and the total count.
func (r *EmployeeRepositoryBun) GetAllEmployeeDeleted(ctx context.Context, page, perPage int) ([]model.Employee, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}
	// count total deleted records
	totalCount, err := tx.NewSelect().Model((*model.Employee)(nil)).WhereDeleted().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	// fetch paginated deleted records
	employees := []model.Employee{}
	err = tx.NewSelect().
		Model(&employees).
		WhereAllWithDeleted().
		Relation("User").
		Relation("User.Address").
		Relation("User.Contact").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return employees, int(totalCount), nil
}

// GetSalaryHistory retorna o histórico salarial do funcionário
func (r *EmployeeRepositoryBun) GetSalaryHistory(ctx context.Context, employeeID uuid.UUID) ([]model.EmployeeSalaryHistory, error) {
	histories := []model.EmployeeSalaryHistory{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&histories).Where("employee_id = ?", employeeID).Order("start_date DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return histories, nil
}

// GetPayments retorna os pagamentos do funcionário
func (r *EmployeeRepositoryBun) GetPayments(ctx context.Context, employeeID uuid.UUID) ([]model.PaymentEmployee, error) {
	payments := []model.PaymentEmployee{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&payments).Where("employee_id = ?", employeeID).Order("payment_date DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *EmployeeRepositoryBun) CreateSalaryHistory(ctx context.Context, h *model.EmployeeSalaryHistory) error {
	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	if _, err := tx.NewInsert().Model(h).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
