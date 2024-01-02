package employeerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type ProductRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProductRepositoryBun(db *bun.DB) *ProductRepositoryBun {
	return &ProductRepositoryBun{db: db}
}

func (r *ProductRepositoryBun) RegisterEmployee(ctx context.Context, p *employeeentity.Employee) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&employeeentity.Employee{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	employee := &employeeentity.Employee{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(employee).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *ProductRepositoryBun) GetAllEmployee(ctx context.Context) ([]employeeentity.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	employees := []employeeentity.Employee{}
	err := r.db.NewSelect().Model(&employees).Scan(ctx)

	return employees, err
}
