package addressdto

import (
	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressDTO struct {
	ID           uuid.UUID                 `json:"id"`
	Street       string                    `json:"street"`
	Number       string                    `json:"number"`
	Complement   string                    `json:"complement"`
	Reference    string                    `json:"reference"`
	Neighborhood string                    `json:"neighborhood"`
	City         string                    `json:"city"`
	UF           string                    `json:"uf"`
	Cep          string                    `json:"cep"`
	AddressType  addressentity.AddressType `json:"address_type"`
	DeliveryTax  float64                   `json:"delivery_tax"`
	Coordinates  Coordinates               `json:"coordinates"`
}

func (a *AddressDTO) FromDomain(address *addressentity.Address) {
	if address == nil {
		return
	}
	coordinates := Coordinates{
		address.Coordinates.Latitude,
		address.Coordinates.Longitude,
	}

	*a = AddressDTO{
		ID:           address.ID,
		Street:       address.Street,
		Number:       address.Number,
		Complement:   address.Complement,
		Reference:    address.Reference,
		Neighborhood: address.Neighborhood,
		City:         address.City,
		UF:           address.UF,
		Cep:          address.Cep,
		AddressType:  address.AddressType,
		DeliveryTax:  address.DeliveryTax,
		Coordinates:  coordinates,
	}
}
