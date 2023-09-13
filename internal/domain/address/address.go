package addressentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	PersonID     uuid.UUID `bun:"person_id,type:uuid,notnull"`
	Street       string    `bun:"street,notnull"`
	Number       string    `bun:"number,notnull"`
	Complement   string    `bun:"complement"`
	Reference    string    `bun:"reference"`
	Neighborhood string    `bun:"neighborhood"`
	City         string    `bun:"city,notnull"`
	State        string    `bun:"state,notnull"`
	Cep          string    `bun:"cep"`
}
