package cliententity

import (
	"time"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	personentity.Person
	TotalOrders  int
	DateRegister time.Time
}
