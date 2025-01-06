package model

import (
	"time"

	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
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

func (c *Coupon) FromDomain(coupon *orderentity.Coupon) {
	*c = Coupon{
		Entity: entitymodel.FromDomain(coupon.Entity),
		CouponCommonAttributes: CouponCommonAttributes{
			Discount: coupon.Discount,
			Min:      coupon.Min,
			StartAt:  coupon.StartAt,
			EndAt:    coupon.EndAt,
		},
	}
}

func (c *Coupon) ToDomain() *orderentity.Coupon {
	return &orderentity.Coupon{
		Entity: c.Entity.ToDomain(),
		CouponCommonAttributes: orderentity.CouponCommonAttributes{
			Discount: c.Discount,
			Min:      c.Min,
			StartAt:  c.StartAt,
			EndAt:    c.EndAt,
		},
	}
}
