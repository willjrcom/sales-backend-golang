package advertisingentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type Advertising struct {
	entity.Entity
	AdvertisingCommonAttributes
}

type AdvertisingCommonAttributes struct {
	Title                      string                 `json:"title"`
	Description                string                 `json:"description"`
	Link                       string                 `json:"link"`
	Contact                    string                 `json:"contact"`
	Type                       string                 `json:"type"`
	StartedAt                  *time.Time             `json:"started_at"`
	EndedAt                    *time.Time             `json:"ended_at"`
	CoverImagePath             string                 `json:"cover_image_path"`
	Images                     []string               `json:"images"`
	SponsorID                  uuid.UUID              `json:"sponsor_id"`
	Sponsor                    *sponsorentity.Sponsor `json:"sponsor"`
	CompanyCategoryAdvertising []CompanyCategory
}

type CompanyCategory struct {
	ID string
}

func NewAdvertising(advertisingCommonAttributes AdvertisingCommonAttributes) *Advertising {
	return &Advertising{
		Entity:                      entity.NewEntity(),
		AdvertisingCommonAttributes: advertisingCommonAttributes,
	}
}
