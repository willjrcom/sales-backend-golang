package employeeusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	re employeeentity.Repository
	rc personentity.ContactRepository
	ru companyentity.UserRepository
}

func NewService(repository employeeentity.Repository) *Service {
	return &Service{re: repository}
}

func (s *Service) AddDependencies(rc personentity.ContactRepository, ru companyentity.UserRepository) {
	s.rc = rc
	s.ru = ru
}

func (s *Service) CreateEmployee(ctx context.Context, dto *employeedto.EmployeeCreateDTO) (*uuid.UUID, error) {
	employee, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	// Get userExists user
	if userExists, _ := s.ru.ExistsUserByID(ctx, *employee.UserID); !userExists {
		return nil, errors.New("user ID not found")
	}

	if employeeExists, _ := s.re.GetEmployeeByUserID(ctx, employee.UserID.String()); employeeExists != nil {
		return nil, errors.New("employee already exists")
	}

	if err := s.re.CreateEmployee(ctx, employee); err != nil {
		return nil, err
	}

	return &employee.ID, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, dtoId *entitydto.IDRequest, dto *employeedto.EmployeeUpdateDTO) error {
	employee, err := s.re.GetEmployeeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateDomain(employee); err != nil {
		return err
	}

	if err := s.re.UpdateEmployee(ctx, employee); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteEmployee(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.re.GetEmployeeById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.re.DeleteEmployee(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetEmployeeById(ctx context.Context, dto *entitydto.IDRequest) (*employeedto.EmployeeDTO, error) {
	if employee, err := s.re.GetEmployeeById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &employeedto.EmployeeDTO{}
		dto.FromDomain(employee)
		return dto, nil
	}
}

func (s *Service) GetEmployeeByUserID(ctx context.Context, dto *entitydto.IDRequest) (*employeedto.EmployeeDTO, error) {
	if employee, err := s.re.GetEmployeeByUserID(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &employeedto.EmployeeDTO{}
		dto.FromDomain(employee)
		return dto, nil
	}
}

func (s *Service) GetAllEmployees(ctx context.Context) ([]employeedto.EmployeeDTO, error) {
	if employees, err := s.re.GetAllEmployees(ctx); err != nil {
		return nil, err
	} else {
		dtos := employeesToDtos(employees)
		return dtos, nil
	}
}

func employeesToDtos(employees []employeeentity.Employee) []employeedto.EmployeeDTO {
	dtos := make([]employeedto.EmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i].FromDomain(&employee)
	}

	return dtos
}
