package employeeusecases

import (
	"context"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type Service struct {
	r employeeentity.Repository
}

func NewService(repository employeeentity.Repository) *Service {
	return &Service{r: repository}
}

func (s *Service) RegisterEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return s.r.RegisterEmployee(ctx, p)
}

func (s *Service) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return s.r.UpdateEmployee(ctx, p)
}

func (s *Service) DeleteEmployee(ctx context.Context, id string) error {
	return s.r.DeleteEmployee(ctx, id)
}

func (s *Service) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	return s.r.GetEmployeeById(ctx, id)
}

func (s *Service) GetAllEmployee(ctx context.Context) ([]employeeentity.Employee, error) {
	return s.r.GetAllEmployee(ctx)
}
