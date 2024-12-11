package deliverydriverusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	deliverydriverdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery_driver"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	r  orderentity.DeliveryDriverRepository
	re employeeentity.Repository
}

func NewService(r orderentity.DeliveryDriverRepository) *Service {
	return &Service{r: r}
}

func (s *Service) AddDependencies(re employeeentity.Repository) {
	s.re = re
}

func (s *Service) CreateDeliveryDriver(ctx context.Context, dto *deliverydriverdto.CreateDeliveryDriverInput) (uuid.UUID, error) {
	driver, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	employee, _ := s.re.GetEmployeeById(ctx, driver.EmployeeID.String())
	if employee == nil {
		return uuid.Nil, errors.New("employee not found")
	}

	oldDeliveryDriver, _ := s.r.GetDeliveryDriverByEmployeeId(ctx, driver.EmployeeID.String())
	if oldDeliveryDriver != nil {
		return uuid.Nil, errors.New("delivery driver already exists")
	}

	if err = s.r.CreateDeliveryDriver(ctx, driver); err != nil {
		return uuid.Nil, err
	}

	return driver.ID, nil
}

func (s *Service) UpdateDeliveryDriver(ctx context.Context, dtoId *entitydto.IdRequest, dto *deliverydriverdto.UpdateDeliveryDriverInput) error {
	driver, err := s.r.GetDeliveryDriverById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(driver); err != nil {
		return err
	}

	if err = s.r.UpdateDeliveryDriver(ctx, driver); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteDeliveryDriver(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetDeliveryDriverById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteDeliveryDriver(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDeliveryDriverByID(ctx context.Context, dto *entitydto.IdRequest) (*deliverydriverdto.DeliveryDriverOutput, error) {
	deliveryDriver, err := s.r.GetDeliveryDriverById(ctx, dto.ID.String())
	if err != nil {
		return nil, err
	}

	deliveryDriverOutput := &deliverydriverdto.DeliveryDriverOutput{}
	deliveryDriverOutput.FromModel(deliveryDriver)
	return deliveryDriverOutput, nil
}

func (s *Service) GetAllDeliveryDrivers(ctx context.Context) ([]deliverydriverdto.DeliveryDriverOutput, error) {
	deliveryDrivers, err := s.r.GetAllDeliveryDrivers(ctx)
	if err != nil {
		return nil, err
	}

	deliveryDriversDto := []deliverydriverdto.DeliveryDriverOutput{}
	for _, deliveryDriver := range deliveryDrivers {
		deliveryDriverOutput := &deliverydriverdto.DeliveryDriverOutput{}
		deliveryDriverOutput.FromModel(&deliveryDriver)
		deliveryDriversDto = append(deliveryDriversDto, *deliveryDriverOutput)
	}

	return deliveryDriversDto, nil
}
