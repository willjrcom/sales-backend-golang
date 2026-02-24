package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CategoryToSponsor struct {
	bun.BaseModel `bun:"table:public.category_sponsors"`
	CategoryID    uuid.UUID `bun:"category_id,type:uuid,pk"`
	SponsorID     uuid.UUID `bun:"sponsor_id,type:uuid,pk"`

	Category *CompanyCategory `bun:"rel:belongs-to,join:category_id=id"`
	Sponsor  *Sponsor         `bun:"rel:belongs-to,join:sponsor_id=id"`
}
