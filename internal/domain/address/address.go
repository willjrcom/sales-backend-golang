package addressentity

import (
	"github.com/google/uuid"
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
	AddressType  AddressType
	DeliveryTax  float64
	Coordinates  Coordinates
}

type AddressType string

const (
	AddressTypeHouse       AddressType = "house"
	AddressTypeApartment   AddressType = "apartment"
	AddressTypeCondominium AddressType = "condominium"
	AddressTypeWork        AddressType = "work"
	AddressTypeHotel       AddressType = "hotel"
	AddressTypeShed        AddressType = "shed"
)

func NewAddress(addressCommonAttributes *AddressCommonAttributes) *Address {
	return &Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: *addressCommonAttributes,
	}
}
