package cliententity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	entity.Entity
	personentity.Person
	TotalOrders int
}
