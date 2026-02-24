package addressentity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	AddressCommonAttributes
	ObjectID uuid.UUID
}

type AddressCommonAttributes struct {
	Street       string
	Number       string
	Complement   string
	Reference    string
	Neighborhood string
	City         string
	UF           string
	Cep          string
	DeliveryTax  decimal.Decimal
	Coordinates  Coordinates
}

func NewAddress(addressCommonAttributes *AddressCommonAttributes) *Address {
	return &Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: *addressCommonAttributes,
	}
}
