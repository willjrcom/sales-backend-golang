package employeeusecases

import (
	"context"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	re employeeentity.Repository
	rc personentity.ContactRepository
}

func NewService(repository employeeentity.Repository) *Service {
	return &Service{re: repository}
}

func (s *Service) AddDependencies(repositoryContact personentity.ContactRepository) {
	s.rc = repositoryContact
}

func (s *Service) CreateEmployee(ctx context.Context, dto *employeedto.CreateEmployeeInput) (uuid.UUID, error) {
	employee, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.re.CreateEmployee(ctx, employee); err != nil {
		return uuid.Nil, err
	}

	return employee.ID, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, dtoId *entitydto.IdRequest, dto *employeedto.UpdateEmployeeInput) error {
	employee, err := s.re.GetEmployeeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(employee); err != nil {
		return err
	}

	if err := s.re.UpdateEmployee(ctx, employee); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteEmployee(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.re.GetEmployeeById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.re.DeleteEmployee(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetEmployeeById(ctx context.Context, dto *entitydto.IdRequest) (*employeedto.EmployeeOutput, error) {
	if employee, err := s.re.GetEmployeeById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &employeedto.EmployeeOutput{}
		dto.FromModel(employee)
		return dto, nil
	}
}

func (s *Service) GetEmployeeByUserID(ctx context.Context, dto *entitydto.IdRequest) (*employeedto.EmployeeOutput, error) {
	if employee, err := s.re.GetEmployeeByUserID(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &employeedto.EmployeeOutput{}
		dto.FromModel(employee)
		return dto, nil
	}
}

func (s *Service) GetAllEmployees(ctx context.Context) ([]employeedto.EmployeeOutput, error) {
	if employees, err := s.re.GetAllEmployees(ctx); err != nil {
		return nil, err
	} else {
		dtos := employeesToDtos(employees)
		return dtos, nil
	}
}

func employeesToDtos(employees []employeeentity.Employee) []employeedto.EmployeeOutput {
	dtos := make([]employeedto.EmployeeOutput, len(employees))
	for i, employee := range employees {
		dtos[i].FromModel(&employee)
	}

	return dtos
}

func (s *Service) CreateContactToEmployee(ctx context.Context, dto *contactdto.CreateContactInput) (uuid.UUID, error) {
	contact, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	// Validate if exists
	if _, err := s.re.GetEmployeeById(ctx, contact.ObjectID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.rc.CreateContact(ctx, contact); err != nil {
		return uuid.Nil, err
	}

	return contact.ID, nil
}

func (s *Service) UpdateContact(ctx context.Context, dtoId *entitydto.IdRequest, dto *contactdto.UpdateContactInput) error {
	contact, err := s.rc.GetContactById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(contact); err != nil {
		return err
	}

	if err := s.rc.UpdateContact(ctx, contact); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteContact(ctx context.Context, dtoId *entitydto.IdRequest) error {
	if _, err := s.rc.GetContactById(ctx, dtoId.ID.String()); err != nil {
		return err
	}

	if err := s.rc.DeleteContact(ctx, dtoId.ID.String()); err != nil {
		return err
	}

	return nil
}
