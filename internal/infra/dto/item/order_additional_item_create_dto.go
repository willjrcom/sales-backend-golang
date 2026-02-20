package itemdto

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type OrderAdditionalItemCreateDTO struct {
	ProductID   uuid.UUID `json:"product_id"`
	VariationID uuid.UUID `json:"variation_id"`
	Quantity    float64   `json:"quantity"`
	Flavor      *string   `json:"flavor,omitempty"`
}

func (a *OrderAdditionalItemCreateDTO) validate() error {
	if a.ProductID == uuid.Nil {
		return errors.New("product id is required")
	}

	if a.VariationID == uuid.Nil {
		return errors.New("variation id is required")
	}

	if a.Quantity == 0 {
		return errors.New("quantity is required")
	}

	return nil
}

func (a *OrderAdditionalItemCreateDTO) ToDomain() (productID uuid.UUID, variationID uuid.UUID, quantity float64, flavor *string, err error) {
	if err = a.validate(); err != nil {
		return
	}

	return a.ProductID, a.VariationID, a.Quantity, a.normalizedFlavor(), nil
}

func (a *OrderAdditionalItemCreateDTO) normalizedFlavor() *string {
	if a.Flavor == nil {
		return nil
	}

	value := strings.TrimSpace(*a.Flavor)
	if value == "" {
		return nil
	}

	return &value
}
