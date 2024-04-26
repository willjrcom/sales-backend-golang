package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrPickupOrderNotReady        = errors.New("pickup order not ready")
	ErrPickupOrderAlreadyReady    = errors.New("pickup order already ready")
	ErrPickupOrderAlreadyPickedup = errors.New("pickup order already pickedup")
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
	ReadyAt    *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
	PickedupAt *time.Time `bun:"pickedup_at" json:"pickedup_at,omitempty"`
}

func NewPickupOrder(name string) *PickupOrder {
	pickupOrderCommonAttributes := PickupOrderCommonAttributes{
		Name:   name,
		Status: PickupOrderStatusPending,
	}

	return &PickupOrder{
		Entity:                      entity.NewEntity(),
		PickupOrderCommonAttributes: pickupOrderCommonAttributes,
	}
}

func (d *PickupOrder) Launch() error {
	if d.Status == PickupOrderStatusReady {
		return ErrPickupOrderAlreadyReady
	}

	if d.Status == PickupOrderStatusPickedup {
		return ErrPickupOrderAlreadyPickedup
	}

	d.ReadyAt = &time.Time{}
	*d.ReadyAt = time.Now()
	d.Status = PickupOrderStatusReady
	return nil
}

func (d *PickupOrder) PickUp() error {
	if d.Status == PickupOrderStatusPickedup {
		return ErrPickupOrderAlreadyPickedup
	}

	if d.Status != PickupOrderStatusReady {
		return ErrPickupOrderNotReady
	}

	d.PickedupAt = &time.Time{}
	*d.PickedupAt = time.Now()
	d.Status = PickupOrderStatusPickedup
	return nil
}
