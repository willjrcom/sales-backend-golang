package model

import (
	"time"

	"github.com/shopspring/decimal"
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
	Discount *decimal.Decimal `bun:"discount,type:decimal(10,2)"`
	Min      *decimal.Decimal `bun:"min,type:decimal(10,2)"`
	StartAt  *time.Time       `bun:"start_at"`
	EndAt    *time.Time       `bun:"end_at"`
}

func (c *Coupon) FromDomain(coupon *orderentity.Coupon) {
	if coupon == nil {
		return
	}
	*c = Coupon{
		Entity: entitymodel.FromDomain(coupon.Entity),
		CouponCommonAttributes: CouponCommonAttributes{
			Discount: &coupon.Discount,
			Min:      &coupon.Min,
			StartAt:  coupon.StartAt,
			EndAt:    coupon.EndAt,
		},
	}
}

func (c *Coupon) ToDomain() *orderentity.Coupon {
	if c == nil {
		return nil
	}
	return &orderentity.Coupon{
		Entity: c.Entity.ToDomain(),
		CouponCommonAttributes: orderentity.CouponCommonAttributes{
			Discount: c.GetDiscount(),
			Min:      c.GetMin(),
			StartAt:  c.StartAt,
			EndAt:    c.EndAt,
		},
	}
}

func (c *Coupon) GetDiscount() decimal.Decimal {
	if c.Discount == nil {
		return decimal.Zero
	}
	return *c.Discount
}

func (c *Coupon) GetMin() decimal.Decimal {
	if c.Min == nil {
		return decimal.Zero
	}
	return *c.Min
}
