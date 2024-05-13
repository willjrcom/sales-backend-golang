package deliverydriverdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameAndActiveIsEmpty = errors.New("name and active can't be empty")
)

type UpdateDeliveryDriverInput struct {
	orderentity.PatchDeliveryDriver
}

func (s *UpdateDeliveryDriverInput) validate() error {
	return nil
}
func (s *UpdateDeliveryDriverInput) UpdateModel(model *orderentity.DeliveryDriver) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	return nil
}
