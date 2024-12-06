package orderentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrOrderPickupMustBePending = errors.New("order pickup must be pending")
)

type OrderPickup struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_pickups,alias:pickup"`
	PickupTimeLogs
	OrderPickupCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type OrderPickupCommonAttributes struct {
	Name    string            `bun:"name,notnull" json:"name"`
	Status  StatusOrderPickup `bun:"status" json:"status"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type PickupTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at" json:"pending_at,omitempty"`
	ReadyAt   *time.Time `bun:"ready_at" json:"ready_at,omitempty"`
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
	*d.PendingAt = time.Now()
	d.Status = OrderPickupStatusPending
	return nil
}

func (d *OrderPickup) Ready() error {
	if d.Status != OrderPickupStatusPending {
		return ErrOrderPickupMustBePending
	}

	d.ReadyAt = &time.Time{}
	*d.ReadyAt = time.Now()
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
