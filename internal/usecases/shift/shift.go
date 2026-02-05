package shiftusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	shiftdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/shift"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
)

var (
	ErrShiftAlreadyClosed = errors.New("shift already closed")
	ErrShiftAlreadyOpened = errors.New("shift already opened")
)

type Service struct {
	r  model.ShiftRepository
	ro model.OrderRepository
	rd model.DeliveryDriverRepository
	se *employeeusecases.Service
	// Repositórios para analytics
	orderProcessRepo model.OrderProcessRepository
	orderQueueRepo   model.QueueRepository
	processRuleRepo  model.ProcessRuleRepository
	employeeRepo     model.EmployeeRepository
}

func NewService(c model.ShiftRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(se *employeeusecases.Service, ro model.OrderRepository, rd model.DeliveryDriverRepository, orderProcessRepo model.OrderProcessRepository,
	orderQueueRepo model.QueueRepository,
	processRuleRepo model.ProcessRuleRepository,
	employeeRepo model.EmployeeRepository) {
	s.se = se
	s.ro = ro
	s.rd = rd
	s.orderProcessRepo = orderProcessRepo
	s.orderQueueRepo = orderQueueRepo
	s.processRuleRepo = processRuleRepo
	s.employeeRepo = employeeRepo
}

func (s *Service) OpenShift(ctx context.Context, dto *shiftdto.ShiftUpdateOpenDTO) (id uuid.UUID, err error) {
	startChange, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	if openedShift, _ := s.r.GetCurrentShift(ctx); openedShift != nil {
		return uuid.Nil, ErrShiftAlreadyOpened
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)

	if !ok {
		return uuid.Nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	employee, err := s.se.GetEmployeeByUserID(ctx, entitydto.NewIdRequest(userIDUUID))
	if err != nil {
		return uuid.Nil, errors.New("user must be an employee")
	}

	shift := shiftentity.NewShift(startChange)
	shift.AttendantID = &employee.ID

	shiftModel := &model.Shift{}
	shiftModel.FromDomain(shift)
	if err = s.r.CreateShift(ctx, shiftModel); err != nil {
		return uuid.Nil, err
	}

	return shift.ID, nil
}

func (s *Service) CloseShift(ctx context.Context, dto *shiftdto.ShiftUpdateCloseDTO) (err error) {
	endChange, err := dto.ToDomain()

	if err != nil {
		return err
	}

	shiftModel, err := s.r.GetFullCurrentShift(ctx)
	if err != nil {
		return err
	}

	shift := shiftModel.ToDomain()
	if shift.IsClosed() {
		return ErrShiftAlreadyClosed
	}

	orderStatus := []orderentity.StatusOrder{
		orderentity.OrderStatusFinished,
		orderentity.OrderStatusCancelled,
	}
	orders, err := s.ro.GetAllOrders(ctx, shiftModel.ID.String(), orderStatus, true, "AND")
	if err != nil {
		return err
	}

	shiftModel.Orders = orders
	shift = shiftModel.ToDomain()

	shift.CloseShift(endChange)

	if err := s.LoadShiftWithProductionAnalytics(ctx, shift); err != nil {
		return err
	}

	shiftModel.FromDomain(shift)
	if err := s.r.UpdateShift(ctx, shiftModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetShiftByID(ctx context.Context, dtoID *entitydto.IDRequest) (shiftDTO *shiftdto.ShiftDTO, err error) {
	shiftModel, err := s.r.GetShiftByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	shift := shiftModel.ToDomain()
	if err := s.LoadShiftWithProductionAnalytics(ctx, shift); err != nil {
		return nil, err
	}

	shiftDTO = &shiftdto.ShiftDTO{}
	shiftDTO.FromDomain(shift)
	return shiftDTO, nil
}

func (s *Service) GetCurrentShift(ctx context.Context) (shiftDTO *shiftdto.ShiftDTO, err error) {
	shiftModel, err := s.r.GetFullCurrentShift(ctx)
	if err != nil {
		return nil, err
	}

	orderStatus := []orderentity.StatusOrder{
		orderentity.OrderStatusFinished,
		orderentity.OrderStatusCancelled,
	}
	orders, err := s.ro.GetAllOrders(ctx, shiftModel.ID.String(), orderStatus, true, "AND")
	if err != nil {
		return nil, err
	}

	shiftModel.Orders = orders

	shift := shiftModel.ToDomain()
	if err := s.LoadShiftWithProductionAnalytics(ctx, shift); err != nil {
		return nil, err
	}

	shiftDTO = &shiftdto.ShiftDTO{}
	shiftDTO.FromDomain(shift)
	return shiftDTO, nil
}

func (s *Service) GetAllShifts(ctx context.Context, page int, perPage int) (shift []shiftdto.ShiftDTO, err error) {
	shiftModels, err := s.r.GetAllShifts(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	shiftDTOs := []shiftdto.ShiftDTO{}
	for _, shiftModel := range shiftModels {
		shift := shiftModel.ToDomain()

		shiftDTO := shiftdto.ShiftDTO{}
		shiftDTO.FromDomain(shift)
		shiftDTOs = append(shiftDTOs, shiftDTO)
	}

	return shiftDTOs, nil
}

func (s *Service) AddRedeem(ctx context.Context, dtoRedeem *shiftdto.ShiftRedeemCreateDTO) (err error) {
	redeem, err := dtoRedeem.ToDomain()
	if err != nil {
		return err
	}

	shiftModel, err := s.r.GetCurrentShift(ctx)
	if err != nil {
		return err
	}

	shift := shiftModel.ToDomain()

	shift.AddRedeem(redeem)

	shiftModel.FromDomain(shift)
	if err := s.r.UpdateShift(ctx, shiftModel); err != nil {
		return err
	}

	return nil
}

func deliveryDriversModelToMap(deliveryDriversModel []model.DeliveryDriver) map[uuid.UUID]orderentity.DeliveryDriver {
	deliveryDrivers := map[uuid.UUID]orderentity.DeliveryDriver{}

	for _, driverModel := range deliveryDriversModel {
		driver := *driverModel.ToDomain()
		deliveryDrivers[driverModel.ID] = driver
	}

	return deliveryDrivers
}

func processRulesModelToMap(processRulesModel []model.ProcessRule) map[uuid.UUID]string {
	processRules := map[uuid.UUID]string{}

	for _, processRuleModel := range processRulesModel {
		processRule := *processRuleModel.ToDomain()
		processRules[processRuleModel.ID] = processRule.Name
	}

	return processRules
}

func employeesModelToMap(employeesModel []model.Employee) map[uuid.UUID]string {
	employees := map[uuid.UUID]string{}

	for _, employeeModel := range employeesModel {
		employee := *employeeModel.ToDomain()
		employees[employeeModel.ID] = employee.User.Name
	}

	return employees
}

// LoadShiftWithProductionAnalytics carrega um shift com todas as métricas de produção
func (s *Service) LoadShiftWithProductionAnalytics(ctx context.Context, shift *shiftentity.Shift) error {
	// Busca todos os processos (vamos usar GetAllProcesses e filtrar por shift depois)
	allProcesses, err := s.orderProcessRepo.GetAllProcesses(ctx)
	if err != nil {
		return err
	}

	// Busca todas as filas
	allQueues, err := s.orderQueueRepo.GetAllQueues(ctx)
	if err != nil {
		return err
	}

	// Busca todas as regras de processo para obter os nomes
	allProcessRules, _, err := s.processRuleRepo.GetAllProcessRules(ctx, 0, 1000, true)
	if err != nil {
		return err
	}

	processRules := processRulesModelToMap(allProcessRules)

	// Busca todos os funcionários para obter os nomes
	allEmployees, _, err := s.employeeRepo.GetAllEmployees(ctx, 0, 1000) // Busca até 1000 funcionários
	if err != nil {
		return err
	}

	employees := employeesModelToMap(allEmployees)

	// Busca delivery drivers para o método Load
	deliveryDriversModel, err := s.rd.GetAllDeliveryDrivers(ctx)
	if err != nil {
		return err
	}
	deliveryDrivers := deliveryDriversModelToMap(deliveryDriversModel)

	// Filtra processos e filas por shift (assumindo que há uma relação com shift)
	// Nota: Esta é uma implementação simplificada. Na prática, você precisaria
	// adicionar campos de shift_id nas tabelas de processos e filas
	var domainProcesses []orderprocessentity.OrderProcess
	for _, p := range allProcesses {
		// Aqui você filtraria por shift_id se existisse
		domainProcesses = append(domainProcesses, *p.ToDomain())
	}

	var domainQueues []*orderprocessentity.OrderQueue
	for _, q := range allQueues {
		// Aqui você filtraria por shift_id se existisse
		domainQueues = append(domainQueues, q.ToDomain())
	}

	// Carrega todas as métricas de uma vez só usando o método Load integrado
	shift.Load(deliveryDrivers, domainProcesses, domainQueues, processRules, employees)

	return nil
}
