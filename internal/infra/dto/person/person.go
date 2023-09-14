package persondto

import (
	"time"

	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type Person struct {
	Name     string               `json:"name"`
	Birthday time.Time            `json:"birthday"`
	Email    string               `json:"email"`
	Contacts []Contact            `json:"contacts"`
	Address  []addressdto.Address `json:"address"`
	Cpf      string               `json:"cpf"`
}
