package deliverydriverdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameAndActiveIsEmpty = errors.New("name and active can't be empty")
)

type DeliveryDriverUpdateDTO struct {
	IsActive *bool `json:"is_active"`
}

func (s *DeliveryDriverUpdateDTO) validate() error {
	return nil
}

func (s *DeliveryDriverUpdateDTO) UpdateDomain(driver *orderentity.DeliveryDriver) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.IsActive != nil {
		driver.IsActive = *s.IsActive
	}

	return nil
}
