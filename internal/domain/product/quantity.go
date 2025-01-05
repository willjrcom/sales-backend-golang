package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrQuantityAlreadyExists = errors.New("quantity already exists")
)

type Quantity struct {
	entity.Entity
	QuantityCommonAttributes
}

type QuantityCommonAttributes struct {
	Quantity   float64
	CategoryID uuid.UUID
}

type PatchQuantity struct {
	Quantity *float64 `json:"quantity"`
}

func NewQuantity(quantityCommonAttributes QuantityCommonAttributes) *Quantity {
	return &Quantity{
		Entity:                   entity.NewEntity(),
		QuantityCommonAttributes: quantityCommonAttributes,
	}
}

func ValidateDuplicateQuantities(name float64, quantities []Quantity) error {
	for _, quantity := range quantities {
		if quantity.Quantity == name {
			return ErrQuantityAlreadyExists
		}
	}

	return nil
}

func ValidateUpdateQuantity(quantity *Quantity, quantities []Quantity) error {
	for _, s := range quantities {
		if s.Quantity == quantity.Quantity && s.ID != quantity.ID {
			return ErrQuantityAlreadyExists
		}
	}

	return nil
}
