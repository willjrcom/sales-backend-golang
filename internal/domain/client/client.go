package cliententity

import (
	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	bun.BaseModel `bun:"table:clients"`
	personentity.Person
	TotalOrders int `bun:"total_orders"`
}
