package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type PublicAddress struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.addresses"`
	AddressCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

func (a *PublicAddress) FromDomain(address *addressentity.Address) {
	if address == nil {
		return
	}
	*a = PublicAddress{
		Entity:   entitymodel.FromDomain(address.Entity),
		ObjectID: address.ObjectID,
		AddressCommonAttributes: AddressCommonAttributes{
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

func (a *PublicAddress) ToDomain() *addressentity.Address {
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
			State:        a.State,
			Cep:          a.Cep,
			AddressType:  addressentity.AddressType(a.AddressType),
			DeliveryTax:  a.DeliveryTax,
		},
	}

	address.Coordinates = *a.Coordinates.ToDomain()

	return address
}
