package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrPickupOrderMustBeStaging = errors.New("pickup order must be staging")
	ErrPickupOrderMustBePending = errors.New("pickup order must be pending")
)

type PickupOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:pickup_orders"`
	PickupTimeLogs
	PickupOrderCommonAttributes
}

type PickupOrderCommonAttributes struct {
	Name    string            `bun:"name,notnull" json:"name"`
	Status  StatusPickupOrder `bun:"status" json:"status"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type PickupTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ReadyAt   *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
}

func NewPickupOrder(name string) *PickupOrder {
	pickupOrderCommonAttributes := PickupOrderCommonAttributes{
		Name:   name,
		Status: PickupOrderStatusStaging,
	}

	return &PickupOrder{
		Entity:                      entity.NewEntity(),
		PickupOrderCommonAttributes: pickupOrderCommonAttributes,
	}
}

func (d *PickupOrder) Pend() error {
	if d.Status != PickupOrderStatusStaging {
		return ErrPickupOrderMustBeStaging
	}

	d.PendingAt = &time.Time{}
	*d.PendingAt = time.Now()
	d.Status = PickupOrderStatusPending
	return nil
}

func (d *PickupOrder) Ready() error {
	if d.Status != PickupOrderStatusPending {
		return ErrPickupOrderMustBePending
	}

	d.ReadyAt = &time.Time{}
	*d.ReadyAt = time.Now()
	d.Status = PickupOrderStatusReady
	return nil
}
