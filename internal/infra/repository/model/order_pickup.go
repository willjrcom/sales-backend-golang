package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

var (
	ErrOrderPickupMustBePending = errors.New("order pickup must be pending")
)

type OrderPickup struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_pickups,alias:pickup"`
	PickupTimeLogs
	OrderPickupCommonAttributes
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
