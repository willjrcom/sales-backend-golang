package ordertabledto

import (
	"time"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderTableDTO struct {
	ID uuid.UUID `json:"id"`
	OrderTableCommonAttributes
	OrderTableTimeLogs
}

type OrderTableCommonAttributes struct {
	Name    string                       `json:"name"`
	Contact string                       `json:"contact"`
	Status  orderentity.StatusOrderTable `json:"status"`
	OrderID uuid.UUID                    `json:"order_id"`
	TableID uuid.UUID                    `json:"table_id"`
}

type OrderTableTimeLogs struct {
	PendingAt *time.Time `json:"pending_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func (t *OrderTableDTO) FromDomain(table *orderentity.OrderTable) {
	if table == nil {
		return
	}
	*t = OrderTableDTO{
		ID: table.ID,
		OrderTableCommonAttributes: OrderTableCommonAttributes{
			Name:    table.Name,
			Contact: table.Contact,
			Status:  table.Status,
			OrderID: table.OrderID,
			TableID: table.TableID,
		},
		OrderTableTimeLogs: OrderTableTimeLogs{
			PendingAt: table.PendingAt,
			ClosedAt:  table.ClosedAt,
		},
	}
}
