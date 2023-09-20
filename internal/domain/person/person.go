package personentity

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Person struct {
	entity.Entity
	Name      string                  `bun:"name,notnull"`
	Email     string                  `bun:"email"`
	Cpf       string                  `bun:"cpf"`
	Birthday  *time.Time              `bun:"birthday"`
	Contacts  []Contact               `bun:"rel:has-many,join:id=person_id,notnull"`
	Addresses []addressentity.Address `bun:"rel:has-many,join:id=person_id,notnull"`
}
