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
	State        string                    `json:"state"`
	Cep          string                    `json:"cep"`
	AddressType  addressentity.AddressType `json:"address_type"`
	DeliveryTax  float64                   `json:"delivery_tax"`
	Coordinates  Coordinates               `json:"coordinates"`
}

func (a *AddressDTO) FromModel(domain *addressentity.Address) {
	coordinates := Coordinates{
		domain.Coordinates.Latitude,
		domain.Coordinates.Longitude,
	}

	*a = AddressDTO{
		ID:           domain.ID,
		Street:       domain.Street,
		Number:       domain.Number,
		Complement:   domain.Complement,
		Reference:    domain.Reference,
		Neighborhood: domain.Neighborhood,
		City:         domain.City,
		State:        domain.State,
		Cep:          domain.Cep,
		AddressType:  domain.AddressType,
		DeliveryTax:  domain.DeliveryTax,
		Coordinates:  coordinates,
	}
}
