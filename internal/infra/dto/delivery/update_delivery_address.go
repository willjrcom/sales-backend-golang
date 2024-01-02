package deliveryorderdto

import (
	"errors"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrInvalidAddressID   = errors.New("address id required")
	ErrAddressNotInClient = errors.New("address not in client")
)

type UpdateDeliveryOrder struct {
	AddressID *uuid.UUID `json:"address_id"`
}

func (u *UpdateDeliveryOrder) Validate() error {
	if u.AddressID == nil {
		return ErrInvalidAddressID
	}
	return nil
}

func (u *UpdateDeliveryOrder) UpdateModel(model *orderentity.DeliveryOrder, address *addressentity.Address) error {
	if err := u.Validate(); err != nil {
		return err
	}

	// Validate if address is from client
	if address.PersonID != model.ClientID {
		return ErrAddressNotInClient
	}

	model.AddressID = *u.AddressID

	return nil
}
