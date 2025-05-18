package employeerepositorybun

import (
	"context"
	"sync"

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

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	// Create employee
	if _, err := r.db.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) UpdateEmployee(ctx context.Context, p *model.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(p).Where("employee.id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	// Delete employee
	if _, err := r.db.NewDelete().Model(&model.Employee{}).Where("employee.id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*model.Employee, error) {
	employee := &model.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.id = ?", id).Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetEmployeeByUserID(ctx context.Context, userID string) (*model.Employee, error) {
	employee := &model.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.user_id = ?", userID).Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	employees := []model.Employee{}
	if err := r.db.NewSelect().Model(&employees).Relation("User").Relation("User.Address").Relation("User.Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepositoryBun) AddPaymentEmployee(ctx context.Context, p *model.PaymentEmployee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}
	if _, err := r.db.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}
	return nil
}
