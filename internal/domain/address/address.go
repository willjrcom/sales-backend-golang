package addressentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	bun.BaseModel `bun:"table:addresses"`
	AddressCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type AddressCommonAttributes struct {
	ObjectID     uuid.UUID   `bun:"object_id,type:uuid,notnull" json:"object_id"`
	Street       string      `bun:"street,notnull" json:"street"`
	Number       string      `bun:"number,notnull" json:"number"`
	Complement   string      `bun:"complement" json:"complement"`
	Reference    string      `bun:"reference" json:"reference"`
	Neighborhood string      `bun:"neighborhood,notnull" json:"neighborhood"`
	City         string      `bun:"city,notnull" json:"city"`
	State        string      `bun:"state,notnull" json:"state"`
	Cep          string      `bun:"cep" json:"cep"`
	AddressType  AddressType `bun:"address_type,notnull" json:"address_type"`
	DeliveryTax  float64     `bun:"delivery_tax,notnull" json:"delivery_tax"`
	Coordinates  Coordinates `bun:"coordinates,type:jsonb" json:"coordinates,omitempty"`
}

type PatchAddress struct {
	Street       *string      `json:"street"`
	Number       *string      `json:"number"`
	Complement   *string      `json:"complement"`
	Reference    *string      `json:"reference"`
	Neighborhood *string      `json:"neighborhood"`
	City         *string      `json:"city"`
	State        *string      `json:"state"`
	Cep          *string      `json:"cep"`
	DeliveryTax  *float64     `json:"delivery_tax"`
	Coordinates  *Coordinates `json:"coordinates,omitempty"`
	AddressType  *AddressType `json:"address_type"`
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

func (a *Address) Validate() error {
	if a.Street == "" {
		return errors.New("street is required")
	}
	if a.Number == "" {
		return errors.New("number is required")
	}
	if a.Neighborhood == "" {
		return errors.New("neighborhood is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	if a.AddressType == "" {
		a.AddressType = AddressTypeHouse
	}

	return nil
}

func NewAddress(addressCommonAttributes *AddressCommonAttributes) *Address {
	coordinates, _ := GetCoordinates(addressCommonAttributes)

	if coordinates != nil {
		addressCommonAttributes.Coordinates = *coordinates
	}

	return &Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: *addressCommonAttributes,
	}
}

func NewPatchAddress(patchAddress *PatchAddress, objectID uuid.UUID) *Address {
	addressCommonAttributes := AddressCommonAttributes{}
	addressCommonAttributes.ObjectID = objectID

	if patchAddress.Cep != nil {
		addressCommonAttributes.Cep = *patchAddress.Cep
	}

	if patchAddress.DeliveryTax != nil {
		addressCommonAttributes.DeliveryTax = *patchAddress.DeliveryTax
	}

	if patchAddress.Street != nil {
		addressCommonAttributes.Street = *patchAddress.Street
	}

	if patchAddress.Number != nil {
		addressCommonAttributes.Number = *patchAddress.Number
	}

	if patchAddress.Complement != nil {
		addressCommonAttributes.Complement = *patchAddress.Complement
	}

	if patchAddress.Reference != nil {
		addressCommonAttributes.Reference = *patchAddress.Reference
	}

	if patchAddress.Neighborhood != nil {
		addressCommonAttributes.Neighborhood = *patchAddress.Neighborhood
	}

	if patchAddress.City != nil {
		addressCommonAttributes.City = *patchAddress.City
	}

	if patchAddress.State != nil {
		addressCommonAttributes.State = *patchAddress.State
	}

	if patchAddress.Coordinates != nil {
		addressCommonAttributes.Coordinates = *patchAddress.Coordinates
	}

	if patchAddress.AddressType != nil {
		addressCommonAttributes.AddressType = *patchAddress.AddressType
	}

	return &Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: addressCommonAttributes,
	}
}
