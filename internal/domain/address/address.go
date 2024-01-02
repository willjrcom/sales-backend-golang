package addressentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	bun.BaseModel `bun:"table:addresses"`
	AddressCommonAttributes
}

type AddressCommonAttributes struct {
	PersonID     uuid.UUID `bun:"person_id,type:uuid,notnull" json:"person_id"`
	Street       string    `bun:"street,notnull" json:"street"`
	Number       string    `bun:"number,notnull" json:"number"`
	Complement   string    `bun:"complement" json:"complement"`
	Reference    string    `bun:"reference" json:"reference"`
	Neighborhood string    `bun:"neighborhood,notnull" json:"neighborhood"`
	City         string    `bun:"city,notnull" json:"city"`
	State        string    `bun:"state,notnull" json:"state"`
	Cep          string    `bun:"cep" json:"cep"`
	IsDefault    bool      `bun:"is_default,notnull" json:"is_default"`
	DeliveryTax  float64   `bun:"delivery_tax,notnull" json:"delivery_tax"`
}

type PatchAddress struct {
	Street       *string  `json:"street"`
	Number       *string  `json:"number"`
	Complement   *string  `json:"complement"`
	Reference    *string  `json:"reference"`
	Neighborhood *string  `json:"neighborhood"`
	City         *string  `json:"city"`
	State        *string  `json:"state"`
	Cep          *string  `json:"cep"`
	DeliveryTax  *float64 `json:"delivery_tax"`
}
