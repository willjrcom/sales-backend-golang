package orderentity

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type TableOrder struct {
	ID       uuid.UUID
	NumTable int
	Waiter   employeeentity.Employee
	QrCode   string
	OrderID  uuid.UUID
}
