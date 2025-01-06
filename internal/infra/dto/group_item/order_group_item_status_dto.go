package groupitemdto

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderGroupItemStatusDTO struct {
	Status orderentity.StatusGroupItem `json:"status"`
}
