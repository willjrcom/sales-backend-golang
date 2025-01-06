package groupitemdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type GroupItemByOrderIDAndStatusDTO struct {
	OrderID uuid.UUID                   `json:"order_id"`
	Status  orderentity.StatusGroupItem `json:"status"`
}
