package advertisingentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Advertising struct {
	entity.Entity
	AdvertisingCommonAttributes
}

type AdvertisingCommonAttributes struct {
	Name      string
	ImagePath string
}

func NewAdvertising(advertisingCommonAttributes AdvertisingCommonAttributes) *Advertising {
	return &Advertising{
		Entity:                      entity.NewEntity(),
		AdvertisingCommonAttributes: advertisingCommonAttributes,
	}
}
