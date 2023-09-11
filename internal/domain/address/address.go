package addressentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	PersonID     uuid.UUID `bun:",notnull"`
	Street       string
	Number       string
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Cep          string
}
