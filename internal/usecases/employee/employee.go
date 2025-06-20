package employeeusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	re model.EmployeeRepository
	rc model.ContactRepository
	ru model.UserRepository
}

func NewService(repository model.EmployeeRepository) *Service {
	return &Service{re: repository}
}

func (s *Service) AddDependencies(rc model.ContactRepository, ru model.UserRepository) {
	s.rc = rc
	s.ru = ru
}

func (s *Service) CreateEmployee(ctx context.Context, dto *employeedto.EmployeeCreateDTO) (*uuid.UUID, error) {
	employee, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	// Get userExists user
	if userExists, _ := s.ru.ExistsUserByID(ctx, employee.UserID); !userExists {
		return nil, errors.New("user ID not found")
	}

	deletedEmployee, _ := s.re.GetEmployeeDeletedByUserID(ctx, employee.UserID.String())
	if deletedEmployee != nil && deletedEmployee.DeletedAt != nil {
		// Reativar funcionário: setar DeletedAt = nil e atualizar campos necessários
		deletedEmployee.DeletedAt = nil

		// Atualize outros campos se necessário
		if err := s.re.UpdateEmployee(ctx, deletedEmployee); err != nil {
			return &deletedEmployee.ID, err
		}

		return &deletedEmployee.ID, nil
	}

	if employeeExists, _ := s.re.GetEmployeeByUserID(ctx, employee.UserID.String()); employeeExists != nil {
		return nil, errors.New("employee already exists")
	}

	employeeModel := &model.Employee{}
	employeeModel.FromDomain(employee)
	if err := s.re.CreateEmployee(ctx, employeeModel); err != nil {
		return nil, err
	}

	return &employee.ID, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, dtoId *entitydto.IDRequest, dto *employeedto.EmployeeUpdateDTO) error {
	employeeModel, err := s.re.GetEmployeeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	employee := employeeModel.ToDomain()
	if err := dto.UpdateDomain(employee); err != nil {
		return err
	}

	employeeModel.FromDomain(employee)
	if err := s.re.UpdateEmployee(ctx, employeeModel); err != nil {
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
	if employeeModel, err := s.re.GetEmployeeById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		employee := employeeModel.ToDomain()
		dto := &employeedto.EmployeeDTO{}
		dto.FromDomain(employee)
		return dto, nil
	}
}

func (s *Service) GetEmployeeByUserID(ctx context.Context, dto *entitydto.IDRequest) (*employeedto.EmployeeDTO, error) {
	if employeeModel, err := s.re.GetEmployeeByUserID(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		employee := employeeModel.ToDomain()
		dto := &employeedto.EmployeeDTO{}
		dto.FromDomain(employee)
		return dto, nil
	}
}

// GetAllEmployees retrieves a paginated list of employees and the total count.
func (s *Service) GetAllEmployees(ctx context.Context, page, perPage int) ([]employeedto.EmployeeDTO, int, error) {
	employeeModels, total, err := s.re.GetAllEmployees(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	dtos := modelsToDTOs(employeeModels)
	return dtos, total, nil
}

// AddPayment records a payment for an employee.
func (s *Service) AddPayment(ctx context.Context, dtoId *entitydto.IDRequest, dtoPayment *employeedto.EmployeePaymentCreateDTO) error {
	if _, err := s.re.GetEmployeeById(ctx, dtoId.ID.String()); err != nil {
		return err
	}

	dtoPayment.EmployeeID = dtoId.ID

	// Buscar histórico salarial vigente
	salaryHistories, err := s.re.GetSalaryHistory(ctx, dtoId.ID)
	if err != nil {
		return err
	}
	var salaryHistoryID *uuid.UUID
	if len(salaryHistories) > 0 {
		salaryHistoryID = &salaryHistories[0].ID
	}

	payment, err := dtoPayment.ToDomain()
	if err != nil {
		return err
	}
	payment.SalaryHistoryID = salaryHistoryID

	paymentModel := &model.PaymentEmployee{}
	paymentModel.FromDomain(payment)
	if err := s.re.AddPaymentEmployee(ctx, paymentModel); err != nil {
		return err
	}

	return nil
}

// GetAllEmployeeDeleted retrieves a paginated list of soft-deleted employees and the total count.
func (s *Service) GetAllEmployeeDeleted(ctx context.Context, page, perPage int) ([]employeedto.EmployeeDTO, int, error) {
	employeeModels, total, err := s.re.GetAllEmployeeDeleted(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	dtos := modelsToDTOs(employeeModels)
	return dtos, total, nil
}

func modelsToDTOs(employeeModels []model.Employee) []employeedto.EmployeeDTO {
	dtos := make([]employeedto.EmployeeDTO, len(employeeModels))
	for i, employeeModel := range employeeModels {
		employee := employeeModel.ToDomain()
		dtos[i].FromDomain(employee)
	}

	return dtos
}

// GetSalaryHistory retorna o histórico salarial do funcionário
func (s *Service) GetSalaryHistory(ctx context.Context, employeeID uuid.UUID) ([]model.EmployeeSalaryHistory, error) {
	return s.re.GetSalaryHistory(ctx, employeeID)
}

// GetPayments retorna os pagamentos do funcionário
func (s *Service) GetPayments(ctx context.Context, employeeID uuid.UUID) ([]model.PaymentEmployee, error) {
	return s.re.GetPayments(ctx, employeeID)
}
