package model

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
	Name    string            `bun:"name,notnull"`
	Status  StatusOrderPickup `bun:"status"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull"`
}

type PickupTimeLogs struct {
	PendingAt *time.Time `bun:"pending_at"`
	ReadyAt   *time.Time `bun:"ready_at"`
}
