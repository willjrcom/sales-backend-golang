package employeeentity

type Repository interface {
	RegisterEmployee(p *Employee) error
	UpdateEmployee(p *Employee) error
	DeleteEmployee(id string) error
	GetEmployeeById(id string) (*Employee, error)
	GetEmployeeBy(key string, value string) (*Employee, error)
	GetAllEmployee(key string, value string) ([]Employee, error)
}
