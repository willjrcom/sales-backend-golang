package orderdeliverydto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrInvalidDriverID = errors.New("driver id required")
)

type UpdateDriverOrder struct {
	DriverID uuid.UUID `json:"driver_id"`
}

func (u *UpdateDriverOrder) validate() error {
	if u.DriverID == uuid.Nil {
		return ErrInvalidDriverID
	}
	return nil
}

func (u *UpdateDriverOrder) UpdateModel(model *orderentity.OrderDelivery) error {
	if err := u.validate(); err != nil {
		return err
	}

	model.DriverID = &u.DriverID

	return nil
}
