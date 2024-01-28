package companyentity

import (
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Company struct {
	entity.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
}

type CompanyCommonAttributes struct {
	SchemaName string                `bun:"schema_name,notnull" json:"schema_name"`
	Name       string                `bun:"name,notnull" json:"name"`
	Cnpj       string                `bun:"cnpj,notnull" json:"cnpj"`
	Email      string                `bun:"email" json:"email"`
	Contacts   []string              `bun:"contacts,type:jsonb" json:"contacts,omitempty"`
	Address    addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull" json:"address,omitempty"`
}
