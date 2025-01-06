package model

import (
	"time"

	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Coupon struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:coupons"`
	CouponCommonAttributes
}

type CouponCommonAttributes struct {
	Discount float64    `bun:"discount"`
	Min      float64    `bun:"min"`
	StartAt  *time.Time `bun:"start_at"`
	EndAt    *time.Time `bun:"end_at"`
}
