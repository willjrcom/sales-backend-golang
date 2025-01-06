package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
	AddressType  AddressType `bun:"address_type,notnull"`
	DeliveryTax  float64     `bun:"delivery_tax,notnull"`
	Coordinates  Coordinates `bun:"coordinates,type:jsonb"`
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
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
