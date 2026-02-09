package employeeusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	re  model.EmployeeRepository
	rc  model.ContactRepository
	ru  model.UserRepository
	rco model.CompanyRepository
}

func NewService(repository model.EmployeeRepository) *Service {
	return &Service{re: repository}
}

func (s *Service) AddDependencies(rc model.ContactRepository, ru model.UserRepository, rco model.CompanyRepository) {
	s.rc = rc
	s.ru = ru
	s.rco = rco
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

	// Se não houver permissões definidas (ou mesmo se houver, queremos garantir defaults?),
	// o usuário pediu "ao criar... ja crie com todas ativadas por padrao".
	// Assumo que se o DTO vier vazio, preenchemos. Se vier preenchido, respeitamos?
	// O DTO normalmente vem do front. Se o front não manda nada, DTO é nil ou vazio.
	// Vamos iterar e setar true para todas que não estiverem no map.
	if employee.Permissions == nil {
		employee.Permissions = make(employeeentity.Permissions)
	}

	allPerms := employeeentity.GetAllPermissions()
	for _, p := range allPerms {
		// Se não existe na lista de input, seta como true
		if _, exists := employee.Permissions[p]; !exists {
			employee.Permissions[p] = true
		}
	}

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

	// Handle active status change
	if dto.IsActive != nil {
		if *dto.IsActive {
			// Activate: Add user to company
			if err := s.rco.AddUserToPublicCompany(ctx, employee.UserID); err != nil {
				return err
			}
		} else {
			// Deactivate: Remove user from company
			if err := s.rco.RemoveUserFromPublicCompany(ctx, employee.UserID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) DeleteEmployee(ctx context.Context, dto *entitydto.IDRequest) error {
	employeeModel, err := s.re.GetEmployeeById(ctx, dto.ID.String())
	if err != nil {
		return err
	}

	if err := s.re.DeleteEmployee(ctx, dto.ID.String()); err != nil {
		return err
	}

	// Remove user from public company
	if err := s.rco.RemoveUserFromPublicCompany(ctx, employeeModel.UserID); err != nil {
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
func (s *Service) GetAllEmployees(ctx context.Context, page, perPage int, isActive bool) ([]employeedto.EmployeeDTO, int, error) {
	employeeModels, total, err := s.re.GetAllEmployees(ctx, page, perPage, isActive)
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

func (s *Service) GetAllEmployeesWithoutDeliveryDrivers(ctx context.Context) ([]employeedto.EmployeeDTO, error) {
	employeeModels, err := s.re.GetAllEmployeesWithoutDeliveryDrivers(ctx)
	if err != nil {
		return nil, err
	}
	dtos := modelsToDTOs(employeeModels)
	return dtos, nil
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
func (s *Service) GetSalaryHistory(ctx context.Context, employeeID uuid.UUID) ([]employeedto.EmployeeSalaryHistoryDTO, error) {
	models, err := s.re.GetSalaryHistory(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	dtos := make([]employeedto.EmployeeSalaryHistoryDTO, len(models))
	for i, m := range models {
		var dto employeedto.EmployeeSalaryHistoryDTO
		dto.FromDomain(m.ToDomain())
		dtos[i] = dto
	}
	return dtos, nil
}

// GetPayments retorna os pagamentos do funcionário
func (s *Service) GetPayments(ctx context.Context, employeeID uuid.UUID) ([]employeedto.EmployeePaymentDTO, error) {
	models, err := s.re.GetPayments(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	dtos := make([]employeedto.EmployeePaymentDTO, len(models))
	for i, m := range models {
		var dto employeedto.EmployeePaymentDTO
		dto.FromDomain(m.ToDomain())
		dtos[i] = dto
	}
	return dtos, nil
}

// Cria um novo histórico salarial para o funcionário
func (s *Service) CreateSalaryHistory(ctx context.Context, employeeID uuid.UUID, dto *employeedto.EmployeeSalaryHistoryCreateDTO) (*employeedto.EmployeeSalaryHistoryDTO, error) {
	if dto.EmployeeID == uuid.Nil {
		dto.EmployeeID = employeeID
	}
	historyDomain := dto.ToDomain()
	historyModel := &model.EmployeeSalaryHistory{}
	historyModel.FromDomain(historyDomain)
	if err := s.re.CreateSalaryHistory(ctx, historyModel); err != nil {
		return nil, err
	}
	historyDTO := &employeedto.EmployeeSalaryHistoryDTO{}
	historyDTO.FromDomain(historyDomain)
	return historyDTO, nil
}

// Cria um novo pagamento para o funcionário
func (s *Service) CreatePayment(ctx context.Context, employeeID uuid.UUID, dto *employeedto.EmployeePaymentCreateDTO) (*employeedto.EmployeePaymentDTO, error) {
	if dto.EmployeeID == uuid.Nil {
		dto.EmployeeID = employeeID
	}
	payment, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}
	// Buscar histórico salarial vigente
	salaryHistories, err := s.re.GetSalaryHistory(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	var salaryHistoryID *uuid.UUID
	if len(salaryHistories) > 0 {
		salaryHistoryID = &salaryHistories[0].ID
	}
	payment.SalaryHistoryID = salaryHistoryID
	paymentModel := &model.PaymentEmployee{}
	paymentModel.FromDomain(payment)
	if err := s.re.AddPaymentEmployee(ctx, paymentModel); err != nil {
		return nil, err
	}
	paymentDTO := &employeedto.EmployeePaymentDTO{}
	paymentDTO.FromDomain(payment)
	return paymentDTO, nil
}
