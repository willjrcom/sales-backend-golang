package orderentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type TableOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:table_orders"`
	OrderID       uuid.UUID                `bun:"column:order_id,type:uuid,notnull"`
	CodTable      string                   `bun:"cod_table,notnull"`
	QrCode        string                   `bun:"qr_code"`
	WaiterID      uuid.UUID                `bun:"column:waiter_id,type:uuid,notnull"`
	Waiter        *employeeentity.Employee `bun:"rel:belongs-to"`
}
