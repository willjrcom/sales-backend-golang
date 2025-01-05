package model

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
	Name     string                 `bun:"name,notnull"`
	CNPJ     string                 `bun:"cnpj,notnull"`
	Email    string                 `bun:"email"`
	Contacts []string               `bun:"contacts,type:jsonb"`
	Address  *addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull"`
}
