package sponsorentity

import (
	"time"

	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Sponsor struct {
	entity.Entity
	bun.BaseModel `bun:"table:sponsors"`
	SponsorCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type SponsorCommonAttributes struct {
	Name     string                 `bun:"name,notnull" json:"name"`
	CNPJ     string                 `bun:"cnpj,notnull" json:"cnpj"`
	Email    string                 `bun:"email" json:"email"`
	Contacts []string               `bun:"contacts,type:jsonb" json:"contacts,omitempty"`
	Address  *addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull" json:"address,omitempty"`
}

func NewSponsor(sponsorCommonAttributes SponsorCommonAttributes) *Sponsor {
	return &Sponsor{
		Entity:                  entity.NewEntity(),
		SponsorCommonAttributes: sponsorCommonAttributes,
	}
}
