package groupitemdto

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type GroupItemByStatusInput struct {
	Status orderentity.StatusGroupItem `json:"status"`
}
