package employeerepositorylocal

import employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"

type ProductRepositoryLocal struct {
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{}
}

func (r *ProductRepositoryLocal) RegisterEmployee(p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryLocal) UpdateEmployee(p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryLocal) DeleteEmployee(id string) error {
	return nil
}

func (r *ProductRepositoryLocal) GetEmployeeById(id string) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetEmployeeBy(key string, value string) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetAllEmployee(key string, value string) ([]employeeentity.Employee, error) {
	return nil, nil
}
