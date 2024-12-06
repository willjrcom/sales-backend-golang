package advertisingentity

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Advertising struct {
	entity.Entity
	bun.BaseModel `bun:"table:advertisements"`
	AdvertisingCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type AdvertisingCommonAttributes struct {
	Name      string `bun:"name,notnull" json:"name"`
	ImagePath string `bun:"image_path" json:"image_path"`
}

type PatchAdvertising struct {
	Name      *string `json:"name"`
	ImagePath *string `json:"image_path"`
}

func NewAdvertising(advertisingCommonAttributes AdvertisingCommonAttributes) *Advertising {
	return &Advertising{
		Entity:                      entity.NewEntity(),
		AdvertisingCommonAttributes: advertisingCommonAttributes,
	}
}
