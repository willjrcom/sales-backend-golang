package employeeusecases

import employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"

type Service struct {
	Repository employeeentity.Repository
}

func NewService(repository employeeentity.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) RegisterEmployee(p *employeeentity.Employee) error {
	return s.Repository.RegisterEmployee(p)
}

func (s *Service) UpdateEmployee(p *employeeentity.Employee) error {
	return s.Repository.UpdateEmployee(p)
}

func (s *Service) DeleteEmployee(id string) error {
	return s.Repository.DeleteEmployee(id)
}

func (s *Service) GetEmployeeById(id string) (*employeeentity.Employee, error) {
	return s.Repository.GetEmployeeById(id)
}

func (s *Service) GetEmployeeBy(key string, value string) (*employeeentity.Employee, error) {
	return s.Repository.GetEmployeeBy(key, value)
}

func (s *Service) GetAllEmployee(key string, value string) ([]employeeentity.Employee, error) {
	return s.Repository.GetAllEmployee(key, value)
}
