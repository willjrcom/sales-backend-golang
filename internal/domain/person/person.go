package personentity

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Person struct {
	entity.Entity
	Name     string
	Birthday time.Time
	Email    string
	Contacts []string
	Address  addressentity.Address
	Cpf      string
}
