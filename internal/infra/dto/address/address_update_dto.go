package addressdto

import (
	"github.com/shopspring/decimal"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressUpdateDTO struct {
	Street       *string                    `json:"street"`
	Number       *string                    `json:"number"`
	Complement   *string                    `json:"complement"`
	Reference    *string                    `json:"reference"`
	Neighborhood *string                    `json:"neighborhood"`
	City         *string                    `json:"city"`
	UF           *string                    `json:"uf"`
	Cep          *string                    `json:"cep"`
	AddressType  *addressentity.AddressType `json:"address_type"`
	DeliveryTax  *decimal.Decimal           `json:"delivery_tax"`
	Coordinates  *Coordinates               `json:"coordinates"`
}

func (a *AddressUpdateDTO) validate() error {
	return nil
}

func (a *AddressUpdateDTO) UpdateDomain(address *addressentity.Address) error {
	if err := a.validate(); err != nil {
		return err
	}

	if a.Cep != nil {
		address.Cep = *a.Cep
	}

	if a.DeliveryTax != nil {
		address.DeliveryTax = *a.DeliveryTax
	}

	if a.Street != nil {
		address.Street = *a.Street
	}

	if a.Number != nil {
		address.Number = *a.Number
	}

	if a.Complement != nil {
		address.Complement = *a.Complement
	}

	if a.Reference != nil {
		address.Reference = *a.Reference
	}

	if a.Neighborhood != nil {
		address.Neighborhood = *a.Neighborhood
	}

	if a.City != nil {
		address.City = *a.City
	}

	if a.UF != nil {
		address.UF = *a.UF
	}

	if a.Coordinates != nil {
		address.Coordinates = addressentity.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}

	if a.AddressType != nil {
		address.AddressType = *a.AddressType
	}

	return nil
}
