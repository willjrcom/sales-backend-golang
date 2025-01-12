package orderdeliverydto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrInvalidDriverID = errors.New("driver id required")
)

type DeliveryOrderDriverUpdateDTO struct {
	DriverID uuid.UUID `json:"driver_id"`
}

func (u *DeliveryOrderDriverUpdateDTO) validate() error {
	if u.DriverID == uuid.Nil {
		return ErrInvalidDriverID
	}
	return nil
}

func (u *DeliveryOrderDriverUpdateDTO) UpdateDomain(delivery *orderentity.OrderDelivery) error {
	if err := u.validate(); err != nil {
		return err
	}

	delivery.DriverID = &u.DriverID

	return nil
}
