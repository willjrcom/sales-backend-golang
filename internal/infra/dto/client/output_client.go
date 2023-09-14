package clientdto

import (
	"time"

	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
)

type Client struct {
	Person       `json:"person"`
	TotalOrders  int       `json:"total_orders"`
	DateRegister time.Time `json:"date_register"`
}

type Person struct {
	Name     string                     `json:"name"`
	Birthday time.Time                  `json:"birthday"`
	Email    string                     `json:"email"`
	Contacts []contactdto.ContactOutput `json:"contacts"`
	Address  []addressdto.Address       `json:"address"`
	Cpf      string                     `json:"cpf"`
}
