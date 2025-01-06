package addressdto

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressUpdateDTO struct {
	Street       *string                    `json:"street"`
	Number       *string                    `json:"number"`
	Complement   *string                    `json:"complement"`
	Reference    *string                    `json:"reference"`
	Neighborhood *string                    `json:"neighborhood"`
	City         *string                    `json:"city"`
	State        *string                    `json:"state"`
	Cep          *string                    `json:"cep"`
	AddressType  *addressentity.AddressType `json:"address_type"`
	DeliveryTax  *float64                   `json:"delivery_tax"`
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
		address.AddressCommonAttributes.Cep = *a.Cep
	}

	if a.DeliveryTax != nil {
		address.AddressCommonAttributes.DeliveryTax = *a.DeliveryTax
	}

	if a.Street != nil {
		address.AddressCommonAttributes.Street = *a.Street
	}

	if a.Number != nil {
		address.AddressCommonAttributes.Number = *a.Number
	}

	if a.Complement != nil {
		address.AddressCommonAttributes.Complement = *a.Complement
	}

	if a.Reference != nil {
		address.AddressCommonAttributes.Reference = *a.Reference
	}

	if a.Neighborhood != nil {
		address.AddressCommonAttributes.Neighborhood = *a.Neighborhood
	}

	if a.City != nil {
		address.AddressCommonAttributes.City = *a.City
	}

	if a.State != nil {
		address.AddressCommonAttributes.State = *a.State
	}

	if a.Coordinates != nil {
		address.AddressCommonAttributes.Coordinates = addressentity.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		}
	}

	if a.AddressType != nil {
		address.AddressCommonAttributes.AddressType = *a.AddressType
	}

	return nil
}
