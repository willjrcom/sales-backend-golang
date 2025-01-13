package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Address struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:addresses"`
	AddressCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

type AddressCommonAttributes struct {
	Street       string      `bun:"street,notnull"`
	Number       string      `bun:"number,notnull"`
	Complement   string      `bun:"complement"`
	Reference    string      `bun:"reference"`
	Neighborhood string      `bun:"neighborhood,notnull"`
	City         string      `bun:"city,notnull"`
	UF           string      `bun:"uf,notnull"`
	Cep          string      `bun:"cep"`
	AddressType  string      `bun:"address_type,notnull"`
	DeliveryTax  float64     `bun:"delivery_tax,notnull"`
	Coordinates  Coordinates `bun:"coordinates,type:jsonb"`
}

func (a *Address) FromDomain(address *addressentity.Address) {
	if address == nil {
		return
	}
	*a = Address{
		Entity:   entitymodel.FromDomain(address.Entity),
		ObjectID: address.ObjectID,
		AddressCommonAttributes: AddressCommonAttributes{
			Street:       address.Street,
			Number:       address.Number,
			Complement:   address.Complement,
			Reference:    address.Reference,
			Neighborhood: address.Neighborhood,
			City:         address.City,
			UF:           address.UF,
			Cep:          address.Cep,
			AddressType:  string(address.AddressType),
			DeliveryTax:  address.DeliveryTax,
		},
	}

	a.Coordinates.FromDomain(&address.Coordinates)
}

func (a *Address) ToDomain() *addressentity.Address {
	if a == nil {
		return nil
	}

	address := &addressentity.Address{
		Entity:   a.Entity.ToDomain(),
		ObjectID: a.ObjectID,
		AddressCommonAttributes: addressentity.AddressCommonAttributes{
			Street:       a.Street,
			Number:       a.Number,
			Complement:   a.Complement,
			Reference:    a.Reference,
			Neighborhood: a.Neighborhood,
			City:         a.City,
			UF:           a.UF,
			Cep:          a.Cep,
			AddressType:  addressentity.AddressType(a.AddressType),
			DeliveryTax:  a.DeliveryTax,
		},
	}

	address.Coordinates = *a.Coordinates.ToDomain()

	return address
}
