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
	Name     string
	CNPJ     string
	Email    string
	Contacts []string
	Address  *addressentity.Address
}

func NewSponsor(sponsorCommonAttributes SponsorCommonAttributes) *Sponsor {
	return &Sponsor{
		Entity:                  entity.NewEntity(),
		SponsorCommonAttributes: sponsorCommonAttributes,
	}
}
