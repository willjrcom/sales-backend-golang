package orderentity

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type TableOrder struct {
	OrderID  uuid.UUID                `bun:"column:order_id,type:uuid,notnull"`
	NumTable int                      `bun:"num_table,notnull"`
	QrCode   string                   `bun:"qr_code"`
	WaiterID uuid.UUID                `bun:"column:waiter_id,type:uuid,notnull"`
	Waiter   *employeeentity.Employee `bun:"rel:belongs-to"`
}
