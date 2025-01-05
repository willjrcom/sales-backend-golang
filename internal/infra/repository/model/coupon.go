package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Coupon struct {
	entity.Entity
	bun.BaseModel `bun:"table:coupons"`
	CouponCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type CouponCommonAttributes struct {
	Discount float64    `bun:"discount"`
	Min      float64    `bun:"min"`
	StartAt  *time.Time `bun:"start_at"`
	EndAt    *time.Time `bun:"end_at"`
}
