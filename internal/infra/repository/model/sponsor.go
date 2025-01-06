package model

import (
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Sponsor struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:sponsors"`
	SponsorCommonAttributes
}

type SponsorCommonAttributes struct {
	Name     string   `bun:"name,notnull"`
	CNPJ     string   `bun:"cnpj,notnull"`
	Email    string   `bun:"email"`
	Contacts []string `bun:"contacts,type:jsonb"`
	Address  *Address `bun:"rel:has-one,join:id=object_id,notnull"`
}
