package addressdto

import (
	"errors"

	"github.com/shopspring/decimal"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

var (
	ErrStreetRequired       = errors.New("street is required")
	ErrNumberRequired       = errors.New("number is required")
	ErrNeighborhoodRequired = errors.New("neighborhood is required")
	ErrCityRequired         = errors.New("city is required")
	ErrUfRequired           = errors.New("uf is required")
	ErrDeliveryTaxRequired  = errors.New("delivery tax is required")
)

type AddressCreateDTO struct {
	Street       string                    `json:"street"`
	Number       string                    `json:"number"`
	Complement   string                    `json:"complement"`
	Reference    string                    `json:"reference"`
	Neighborhood string                    `json:"neighborhood"`
	City         string                    `json:"city"`
	UF           string                    `json:"uf"`
	Cep          string                    `json:"cep"`
	AddressType  addressentity.AddressType `json:"address_type"`
	DeliveryTax  *decimal.Decimal          `json:"delivery_tax"`
	Coordinates  Coordinates               `json:"coordinates"`
}

func (a *AddressCreateDTO) validate(withDeliveryTax bool) error {
	if a.Street == "" {
		return ErrStreetRequired
	}
	if a.Number == "" {
		return ErrNumberRequired
	}
	if a.Neighborhood == "" {
		return ErrNeighborhoodRequired
	}
	if a.City == "" {
		return ErrCityRequired
	}
	if a.UF == "" {
		return ErrUfRequired
	}
	if a.AddressType == "" {
		house := addressentity.AddressTypeHouse
		a.AddressType = house
	}

	if withDeliveryTax && a.DeliveryTax == nil {
		return ErrDeliveryTaxRequired
	}

	return nil
}

func (a *AddressCreateDTO) ToDomain(withDeliveryTax bool) (*addressentity.Address, error) {
	if err := a.validate(withDeliveryTax); err != nil {
		return nil, err
	}

	addressCommonAttributes := addressentity.AddressCommonAttributes{
		Street:       a.Street,
		Number:       a.Number,
		Complement:   a.Complement,
		Reference:    a.Reference,
		Neighborhood: a.Neighborhood,
		City:         a.City,
		UF:           a.UF,
		Cep:          a.Cep,
		AddressType:  a.AddressType,
	}

	if withDeliveryTax {
		addressCommonAttributes.DeliveryTax = *a.DeliveryTax
	}

	return addressentity.NewAddress(&addressCommonAttributes), nil
}
