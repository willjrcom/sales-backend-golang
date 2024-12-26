package employeerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
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

func (r *EmployeeRepositoryBun) CreateEmployee(ctx context.Context, c *employeeentity.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Create employee
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if c.User.Person.Contact != nil {
		if _, err := tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", c.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(c.User.Person.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if c.User.Person.Address != nil {
		if _, err := tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", c.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(c.User.Person.Address).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(p).Where("employee.id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	if p.User.Person.Contact != nil {
		if _, err := tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", p.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(p.User.Person.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if p.User.Person.Address != nil {
		if _, err := tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", p.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(p.User.Person.Address).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
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

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Delete employee
	if _, err = tx.NewDelete().Model(&employeeentity.Employee{}).Where("employee.id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete contact
	if _, err = tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete address
	if _, err = tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.id = ?", id).Relation("Address").Relation("Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetEmployeeByUserID(ctx context.Context, userID string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(employee).Where("employee.user_id = ?", userID).Relation("Address").Relation("Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryBun) GetAllEmployees(ctx context.Context) ([]employeeentity.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	employees := []employeeentity.Employee{}
	if err := r.db.NewSelect().Model(&employees).Relation("Address").Relation("Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return employees, nil
}
