package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Product struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:products"`
	ProductCommonAttributes
}

type ProductCommonAttributes struct {
	Code        string           `bun:"code,notnull"`
	Name        string           `bun:"name,notnull"`
	Flavors     []string         `bun:"flavors,type:jsonb"`
	ImagePath   *string          `bun:"image_path"`
	Description string           `bun:"description"`
	Price       float64          `bun:"price,notnull"`
	Cost        float64          `bun:"cost"`
	IsAvailable bool             `bun:"is_available"`
	CategoryID  uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category    *ProductCategory `bun:"rel:belongs-to"`
	SizeID      uuid.UUID        `bun:"column:size_id,type:uuid,notnull"`
	Size        *Size            `bun:"rel:belongs-to"`
}
