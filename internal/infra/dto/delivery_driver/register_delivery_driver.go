package deliverydriverdto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrEmployeeIDRequired = errors.New("employee ID is required")
)

type CreateDeliveryDriverInput struct {
	orderentity.DeliveryDriverCommonAttributes
}

func (s *CreateDeliveryDriverInput) validate() error {
	if s.EmployeeID == uuid.Nil {
		return ErrEmployeeIDRequired
	}

	return nil
}

func (s *CreateDeliveryDriverInput) ToModel() (*orderentity.DeliveryDriver, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	deliveryDriverCommonAttributes := orderentity.DeliveryDriverCommonAttributes{
		EmployeeID: s.EmployeeID,
	}

	return orderentity.NewDeliveryDriver(deliveryDriverCommonAttributes), nil
}
