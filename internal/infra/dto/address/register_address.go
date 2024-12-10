package addressdto

import (
	"errors"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrStreetRequired       = errors.New("street is required")
	ErrNumberRequired       = errors.New("number is required")
	ErrNeighborhoodRequired = errors.New("neighborhood is required")
	ErrCityRequired         = errors.New("city is required")
	ErrStateRequired        = errors.New("state is required")
)

type CreateAddressInput struct {
	addressentity.PatchAddress
}

func (a *CreateAddressInput) validate() error {
	if a.Street == nil {
		return ErrStreetRequired
	}
	if a.Number == nil {
		return ErrNumberRequired
	}
	if a.Neighborhood == nil {
		return ErrNeighborhoodRequired
	}
	if a.City == nil {
		return ErrCityRequired
	}
	if a.State == nil {
		return ErrStateRequired
	}
	if a.AddressType == nil {
		house := addressentity.AddressTypeHouse
		a.AddressType = &house
	}

	return nil
}

func (a *CreateAddressInput) ToModel() (*addressentity.Address, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	addressCommonAttributes := addressentity.AddressCommonAttributes{
		Street:       *a.Street,
		Number:       *a.Number,
		Complement:   *a.Complement,
		Reference:    *a.Reference,
		Neighborhood: *a.Neighborhood,
		City:         *a.City,
		State:        *a.State,
		Cep:          *a.Cep,
		AddressType:  *a.AddressType,
	}

	return &addressentity.Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: addressCommonAttributes,
	}, nil
}
