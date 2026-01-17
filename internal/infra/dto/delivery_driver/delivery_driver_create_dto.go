package deliverydriverdto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrEmployeeIDRequired = errors.New("employee ID is required")
)

type DeliveryDriverCreateDTO struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	IsActive   *bool     `json:"is_active"`
}

func (s *DeliveryDriverCreateDTO) validate() error {
	if s.EmployeeID == uuid.Nil {
		return ErrEmployeeIDRequired
	}

	return nil
}

func (s *DeliveryDriverCreateDTO) ToDomain() (*orderentity.DeliveryDriver, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	isActive := true
	if s.IsActive != nil {
		isActive = *s.IsActive
	}

	deliveryDriverCommonAttributes := orderentity.DeliveryDriverCommonAttributes{
		EmployeeID: s.EmployeeID,
		IsActive:   isActive,
	}

	return orderentity.NewDeliveryDriver(deliveryDriverCommonAttributes), nil
}
