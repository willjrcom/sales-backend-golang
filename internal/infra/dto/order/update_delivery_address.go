package orderdto

import (
	"errors"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrInvalidAddressID   = errors.New("address id required")
	ErrDeliveryNotFound   = errors.New("delivery not found")
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

func (u *UpdateDeliveryOrder) UpdateModel(model *orderentity.Order, address *addressentity.Address) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if model.Delivery == nil {
		return ErrDeliveryNotFound
	}

	// Invalid person
	if address.PersonID != (*model.Delivery).ClientID {
		return ErrAddressNotInClient
	}

	(*model.Delivery).AddressID = *u.AddressID

	return nil
}
