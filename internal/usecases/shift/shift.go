package shiftusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
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
	se *employeeusecases.Service
}

func NewService(c model.ShiftRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(se *employeeusecases.Service, ro model.OrderRepository) {
	s.se = se
	s.ro = ro
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
		orderentity.OrderStatusCanceled,
	}
	orders, err := s.ro.GetAllOrders(ctx, shiftModel.ID.String(), orderStatus, true)
	if err != nil {
		return err
	}

	shiftModel.Orders = orders
	shift = shiftModel.ToDomain()

	shift.CloseShift(endChange)
	shift.Load()

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
	shift.Load()
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
		orderentity.OrderStatusCanceled,
	}
	orders, err := s.ro.GetAllOrders(ctx, shiftModel.ID.String(), orderStatus, true)
	if err != nil {
		return nil, err
	}

	shiftModel.Orders = orders

	shift := shiftModel.ToDomain()
	shift.Load()

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
