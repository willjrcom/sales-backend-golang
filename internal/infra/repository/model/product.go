package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("product category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	bun.BaseModel `bun:"table:products"`
	ProductCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
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
