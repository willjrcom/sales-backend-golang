package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrSizeAlreadyExists = errors.New("size already exists")
)

type Size struct {
	entity.Entity
	bun.BaseModel `bun:"table:sizes"`
	SizeCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type SizeCommonAttributes struct {
	Name       string    `bun:"name"`
	IsActive   *bool     `bun:"is_active"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	Products   []Product `bun:"rel:has-many,join:id=size_id"`
}
