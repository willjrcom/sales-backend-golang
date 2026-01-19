package employeerepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type EmployeeRepositoryBun struct {
	db *bun.DB
}

func NewEmployeeRepositoryBun(db *bun.DB) model.EmployeeRepository {
	return &EmployeeRepositoryBun{db: db}
}

func (r *EmployeeRepositoryBun) CreateEmployee(ctx context.Context, c *model.Employee) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(p).Where("employee.id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Delete employee
	if _, err := tx.NewUpdate().Model(&model.Employee{}).Set("is_active = false").Where("employee.id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*model.Employee, error) {
	employee := &model.Employee{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(employee).Where("employee.user_id = ?", userID).WhereAllWithDeleted().Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employee, nil
}

// GetAllEmployees retrieves a paginated list of employees and the total count.
func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context, page, perPage int, isActive ...bool) ([]model.Employee, int, error) {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	activeFilter := true
	if len(isActive) > 0 {
		activeFilter = isActive[0]
	}

	// count total records
	totalCount, err := tx.NewSelect().Model((*model.Employee)(nil)).Where("employee.is_active = ?", activeFilter).Count(ctx)
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
		Where("employee.is_active = ?", activeFilter).
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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()
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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

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

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&payments).Where("employee_id = ?", employeeID).Order("payment_date DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *EmployeeRepositoryBun) CreateSalaryHistory(ctx context.Context, h *model.EmployeeSalaryHistory) error {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()
	if _, err := tx.NewInsert().Model(h).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetAllEmployeesWithoutDeliveryDrivers retrieves all employees who are not delivery drivers.
func (r *EmployeeRepositoryBun) GetAllEmployeesWithoutDeliveryDrivers(ctx context.Context) ([]model.Employee, error) {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	employees := []model.Employee{}
	err = tx.NewSelect().
		Model(&employees).
		Relation("User").
		Relation("User.Address").
		Relation("User.Contact").
		Where("NOT EXISTS (SELECT 1 FROM delivery_drivers AS dd WHERE dd.employee_id = employee.id)").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return employees, nil
}
