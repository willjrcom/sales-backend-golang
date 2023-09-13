package personentity

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Person struct {
	entity.Entity
	Name     string                  `bun:"name,notnull"`
	Birthday time.Time               `bun:"birthday"`
	Email    string                  `bun:"email"`
	Contacts []string                `bun:"contacts,notnull"`
	Address  []addressentity.Address `bun:"rel:has-many,join:id=person_id,notnull"`
	Cpf      string                  `bun:"cpf"`
}
