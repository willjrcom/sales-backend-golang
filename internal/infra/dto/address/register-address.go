package addressdto

import (
	"errors"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrStreetRequired = errors.New("street is required")
	ErrNumberRequired = errors.New("number is required")
	ErrStateRequired  = errors.New("state is required")
	ErrCityRequired   = errors.New("city is required")
)

type RegisterAddressInput struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Reference    string `json:"reference"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Cep          string `json:"cep"`
}

func (a *RegisterAddressInput) validate() error {
	if a.Street == "" {
		return ErrStreetRequired
	}
	if a.Number == "" {
		return ErrNumberRequired
	}
	if a.State == "" {
		return ErrStateRequired
	}
	if a.City == "" {
		return ErrCityRequired
	}
	return nil
}

func (a *RegisterAddressInput) ToModel() (*addressentity.Address, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return &addressentity.Address{
		Entity:       entity.NewEntity(),
		Street:       a.Street,
		Number:       a.Number,
		Complement:   a.Complement,
		Reference:    a.Reference,
		Neighborhood: a.Neighborhood,
		City:         a.City,
		State:        a.State,
		Cep:          a.Cep,
	}, nil
}
