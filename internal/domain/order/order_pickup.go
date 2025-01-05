package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderPickupMustBePending = errors.New("order pickup must be pending")
)

type OrderPickup struct {
	entity.Entity
	PickupTimeLogs
	OrderPickupCommonAttributes
}

type OrderPickupCommonAttributes struct {
	Name    string
	Status  StatusOrderPickup
	OrderID uuid.UUID
}

type PickupTimeLogs struct {
	PendingAt *time.Time
	ReadyAt   *time.Time
}

func NewOrderPickup(name string) *OrderPickup {
	orderPickupCommonAttributes := OrderPickupCommonAttributes{
		Name:   name,
		Status: OrderPickupStatusStaging,
	}

	return &OrderPickup{
		Entity:                      entity.NewEntity(),
		OrderPickupCommonAttributes: orderPickupCommonAttributes,
	}
}

func (d *OrderPickup) Pend() error {
	if d.Status != OrderPickupStatusStaging {
		return nil
	}

	d.PendingAt = &time.Time{}
	*d.PendingAt = time.Now().UTC()
	d.Status = OrderPickupStatusPending
	return nil
}

func (d *OrderPickup) Ready() error {
	if d.Status != OrderPickupStatusPending {
		return ErrOrderPickupMustBePending
	}

	d.ReadyAt = &time.Time{}
	*d.ReadyAt = time.Now().UTC()
	d.Status = OrderPickupStatusReady
	return nil
}

func (d *OrderPickup) UpdateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	d.Name = name
	return nil
}
