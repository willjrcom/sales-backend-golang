package ordertabledto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderTableDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	OrderTableCommonAttributes
	OrderTableTimeLogs
}

type OrderTableCommonAttributes struct {
	Name    string                       `json:"name"`
	Contact string                       `json:"contact"`
	Status  orderentity.StatusOrderTable `json:"status"`
	TaxRate decimal.Decimal              `json:"tax_rate"`
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
		ID:        table.ID,
		CreatedAt: table.CreatedAt,
		OrderTableCommonAttributes: OrderTableCommonAttributes{
			Name:    table.Name,
			Contact: table.Contact,
			Status:  table.Status,
			TaxRate: table.TaxRate,
			OrderID: table.OrderID,
			TableID: table.TableID,
		},
		OrderTableTimeLogs: OrderTableTimeLogs{
			PendingAt: table.PendingAt,
			ClosedAt:  table.ClosedAt,
		},
	}
}
