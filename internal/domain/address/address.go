package addressentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	Street       string
	Number       string
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	State        string
	Cep          string
}
