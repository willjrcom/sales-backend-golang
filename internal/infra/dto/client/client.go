package clientdto

import (
	"time"

	persondto "github.com/willjrcom/sales-backend-go/internal/infra/dto/person"
)

type Client struct {
	persondto.Person `json:"person"`
	TotalOrders      int       `json:"total_orders"`
	DateRegister     time.Time `json:"date_register"`
}
