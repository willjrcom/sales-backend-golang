package groupitemdto

import groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"

type GroupStatusInput struct {
	Status groupitementity.StatusGroupItem `json:"status"`
}
