package sponsorentity

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Sponsor struct {
	entity.Entity
	SponsorCommonAttributes
}

type SponsorCommonAttributes struct {
	Name                   string
	CNPJ                   string
	Email                  string
	Contact                string
	Address                *addressentity.Address
	CompanyCategorySponsor []CompanyCategory
}

type CompanyCategory struct {
	ID string
}

func NewSponsor(sponsorCommonAttributes SponsorCommonAttributes) *Sponsor {
	return &Sponsor{
		Entity:                  entity.NewEntity(),
		SponsorCommonAttributes: sponsorCommonAttributes,
	}
}
