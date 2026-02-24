package addressdto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressDTO struct {
	ID           uuid.UUID       `json:"id"`
	Street       string          `json:"street"`
	Number       string          `json:"number"`
	Complement   string          `json:"complement"`
	Reference    string          `json:"reference"`
	Neighborhood string          `json:"neighborhood"`
	City         string          `json:"city"`
	UF           string          `json:"uf"`
	Cep          string          `json:"cep"`
	DeliveryTax  decimal.Decimal `json:"delivery_tax"`
	Coordinates  Coordinates     `json:"coordinates"`
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
		DeliveryTax:  address.DeliveryTax,
		Coordinates:  coordinates,
	}
}

func (a *AddressDTO) ToDomain() (*addressentity.Address, error) {
	if a == nil {
		return nil, nil
	}
	coordinates := addressentity.Coordinates{
		Latitude:  a.Coordinates.Latitude,
		Longitude: a.Coordinates.Longitude,
	}

	return addressentity.NewAddress(&addressentity.AddressCommonAttributes{
		Street:       a.Street,
		Number:       a.Number,
		Complement:   a.Complement,
		Reference:    a.Reference,
		Neighborhood: a.Neighborhood,
		City:         a.City,
		UF:           a.UF,
		Cep:          a.Cep,
		DeliveryTax:  a.DeliveryTax,
		Coordinates:  coordinates,
	}), nil
}

func (a *AddressDTO) UpdateDomain(address *addressentity.Address) error {
	if a == nil || address == nil {
		return nil
	}
	address.Street = a.Street
	address.Number = a.Number
	address.Complement = a.Complement
	address.Reference = a.Reference
	address.Neighborhood = a.Neighborhood
	address.City = a.City
	address.UF = a.UF
	address.Cep = a.Cep
	address.DeliveryTax = a.DeliveryTax
	address.Coordinates.Latitude = a.Coordinates.Latitude
	address.Coordinates.Longitude = a.Coordinates.Longitude
	return nil
}
