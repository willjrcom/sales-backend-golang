package groupitemdto

import (
	"github.com/google/uuid"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
)

type GroupItemByOrderIDAndStatusInput struct {
	OrderID uuid.UUID                       `json:"order_id"`
	Status  groupitementity.StatusGroupItem `json:"status"`
}
