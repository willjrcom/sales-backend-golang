package orderentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type Order struct {
	entity.Entity
	bun.BaseModel `bun:"table:orders"`
	PaymentOrder
	Name        string                   `bun:"name,notnull"`
	OrderNumber int                      `bun:"order_number,notnull"`
	Delivery    *DeliveryOrder           `bun:"rel:has-one,join:id=order_id"`
	TableOrder  *TableOrder              `bun:"rel:has-one,join:id=order_id"`
	AttendantID uuid.UUID                `bun:"column:attendant_id,type:uuid,notnull"`
	Attendant   *employeeentity.Employee `bun:"rel:belongs-to"`
	Status      StatusOrder              `bun:"status,type:enum,notnull"`
	Observation string                   `bun:"observation"`
	GroupItems  []itementity.GroupItem   `bun:"rel:has-many,join:id=order_id"`
}
