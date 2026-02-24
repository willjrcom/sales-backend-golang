package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CategoryToAdvertising struct {
	bun.BaseModel `bun:"table:public.category_advertisements"`
	CategoryID    uuid.UUID `bun:"category_id,type:uuid,pk"`
	AdvertisingID uuid.UUID `bun:"advertising_id,type:uuid,pk"`

	Category    *CompanyCategory `bun:"rel:belongs-to,join:category_id=id"`
	Advertising *Advertising     `bun:"rel:belongs-to,join:advertising_id=id"`
}
