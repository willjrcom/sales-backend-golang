package deliverydriverusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	deliverydriverdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery_driver"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	r  model.DeliveryDriverRepository
	re model.EmployeeRepository
}

func NewService(r model.DeliveryDriverRepository) *Service {
	return &Service{r: r}
}

func (s *Service) AddDependencies(re model.EmployeeRepository) {
	s.re = re
}

func (s *Service) CreateDeliveryDriver(ctx context.Context, dto *deliverydriverdto.DeliveryDriverCreateDTO) (uuid.UUID, error) {
	driver, err := dto.ToDomain()

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

func (s *Service) UpdateDeliveryDriver(ctx context.Context, dtoId *entitydto.IDRequest, dto *deliverydriverdto.DeliveryDriverUpdateDTO) error {
	driver, err := s.r.GetDeliveryDriverById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateDomain(driver); err != nil {
		return err
	}

	if err = s.r.UpdateDeliveryDriver(ctx, driver); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteDeliveryDriver(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.r.GetDeliveryDriverById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteDeliveryDriver(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDeliveryDriverByID(ctx context.Context, dto *entitydto.IDRequest) (*deliverydriverdto.DeliveryDriverDTO, error) {
	deliveryDriver, err := s.r.GetDeliveryDriverById(ctx, dto.ID.String())
	if err != nil {
		return nil, err
	}

	deliveryDriverDTO := &deliverydriverdto.DeliveryDriverDTO{}
	deliveryDriverDTO.FromDomain(deliveryDriver)
	return deliveryDriverDTO, nil
}

func (s *Service) GetAllDeliveryDrivers(ctx context.Context) ([]deliverydriverdto.DeliveryDriverDTO, error) {
	deliveryDrivers, err := s.r.GetAllDeliveryDrivers(ctx)
	if err != nil {
		return nil, err
	}

	deliveryDriversDto := []deliverydriverdto.DeliveryDriverDTO{}
	for _, deliveryDriver := range deliveryDrivers {
		deliveryDriverDTO := &deliverydriverdto.DeliveryDriverDTO{}
		deliveryDriverDTO.FromDomain(&deliveryDriver)
		deliveryDriversDto = append(deliveryDriversDto, *deliveryDriverDTO)
	}

	return deliveryDriversDto, nil
}
