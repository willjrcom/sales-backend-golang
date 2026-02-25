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
	Street       string           `json:"street"`
	Number       string           `json:"number"`
	Complement   string           `json:"complement"`
	Reference    string           `json:"reference"`
	Neighborhood string           `json:"neighborhood"`
	City         string           `json:"city"`
	UF           string           `json:"uf"`
	Cep          string           `json:"cep"`
	DeliveryTax  *decimal.Decimal `json:"delivery_tax"`
	Distance     *float64         `json:"distance"`
	Coordinates  Coordinates      `json:"coordinates"`
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
	}

	if withDeliveryTax {
		addressCommonAttributes.DeliveryTax = *a.DeliveryTax
	}

	if a.Distance != nil {
		addressCommonAttributes.Distance = *a.Distance
	}

	return addressentity.NewAddress(&addressCommonAttributes), nil
}
