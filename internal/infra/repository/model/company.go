package model

import (
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Company struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
}

type CompanyCommonAttributes struct {
	SchemaName   string   `bun:"schema_name,notnull"`
	BusinessName string   `bun:"business_name,notnull"`
	TradeName    string   `bun:"trade_name,notnull"`
	Cnpj         string   `bun:"cnpj,notnull"`
	Email        string   `bun:"email"`
	Contacts     []string `bun:"contacts,type:jsonb"`
	Address      *Address `bun:"rel:has-one,join:id=object_id,notnull"`
}
