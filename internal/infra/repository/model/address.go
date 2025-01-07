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
}

type AddressCommonAttributes struct {
	ObjectID     uuid.UUID   `bun:"object_id,type:uuid,notnull"`
	Street       string      `bun:"street,notnull"`
	Number       string      `bun:"number,notnull"`
	Complement   string      `bun:"complement"`
	Reference    string      `bun:"reference"`
	Neighborhood string      `bun:"neighborhood,notnull"`
	City         string      `bun:"city,notnull"`
	State        string      `bun:"state,notnull"`
	Cep          string      `bun:"cep"`
	AddressType  string      `bun:"address_type,notnull"`
	DeliveryTax  float64     `bun:"delivery_tax,notnull"`
	Coordinates  Coordinates `bun:"coordinates,type:jsonb"`
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func (c *Coordinates) FromDomain(coordinates *addressentity.Coordinates) {
	if coordinates == nil {
		return
	}
	*c = Coordinates{
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
	}
}

func (c *Coordinates) ToDomain() *addressentity.Coordinates {
	return &addressentity.Coordinates{
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}
}

func (a *Address) FromDomain(address *addressentity.Address) {
	if address == nil {
		return
	}
	*a = Address{
		Entity: entitymodel.FromDomain(address.Entity),
		AddressCommonAttributes: AddressCommonAttributes{
			ObjectID:     address.ObjectID,
			Street:       address.Street,
			Number:       address.Number,
			Complement:   address.Complement,
			Reference:    address.Reference,
			Neighborhood: address.Neighborhood,
			City:         address.City,
			State:        address.State,
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
		Entity: a.Entity.ToDomain(),
		AddressCommonAttributes: addressentity.AddressCommonAttributes{
			ObjectID:     a.ObjectID,
			Street:       a.Street,
			Number:       a.Number,
			Complement:   a.Complement,
			Reference:    a.Reference,
			Neighborhood: a.Neighborhood,
			City:         a.City,
			State:        a.State,
			Cep:          a.Cep,
			AddressType:  addressentity.AddressType(a.AddressType),
			DeliveryTax:  a.DeliveryTax,
		},
	}

	address.Coordinates = *a.Coordinates.ToDomain()

	return address
}
