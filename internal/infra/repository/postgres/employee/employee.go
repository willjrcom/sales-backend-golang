package employeerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type EmployeeRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewEmployeeRepositoryBun(db *bun.DB) *EmployeeRepositoryBun {
	return &EmployeeRepositoryBun{db: db}
}

func (r *EmployeeRepositoryBun) RegisterEmployee(ctx context.Context, c *employeeentity.Employee) error {
	r.mu.Lock()
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Register employee
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// Register contact
	if _, err := tx.NewInsert().Model(&c.Contact).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Register address
	if _, err := tx.NewInsert().Model(&c.Address).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func rollback(tx *bun.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return err
}

func (r *EmployeeRepositoryBun) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	r.mu.Lock()
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Delete employee
	if _, err = tx.NewDelete().Model(&employeeentity.Employee{}).Where("id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Delete contact
	if _, err = tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Delete address
	if _, err = tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(employee).Where("id = ?", id).Relation("Address").Relation("Contact").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context) ([]employeeentity.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	employees := []employeeentity.Employee{}
	err := r.db.NewSelect().Model(&employees).Relation("Address").Relation("Contact").Scan(ctx)

	return employees, err
}
