package orderentity

import (
	"errors"
	"time"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrDiscountMustBePositive = errors.New("discount must be positive")
	ErrStartAndEndAtRequired  = errors.New("start_at and end_at are required")
	ErrStartAtAfterEndAt      = errors.New("start_at must be before end_at")
)

type Coupon struct {
	entity.Entity
	bun.BaseModel `bun:"table:coupons"`
	CouponCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type CouponCommonAttributes struct {
	Discount float64    `bun:"discount" json:"discount"`
	Min      float64    `bun:"min" json:"min"`
	StartAt  *time.Time `bun:"start_at" json:"start_at"`
	EndAt    *time.Time `bun:"end_at" json:"end_at"`
}

func NewCoupon(couponCommonAttributes CouponCommonAttributes) (*Coupon, error) {
	if couponCommonAttributes.Discount <= 0 {
		return nil, ErrDiscountMustBePositive
	}

	if couponCommonAttributes.StartAt == nil || couponCommonAttributes.EndAt == nil {
		return nil, ErrStartAndEndAtRequired
	}

	if (*couponCommonAttributes.StartAt).After(*couponCommonAttributes.EndAt) {
		return nil, ErrStartAtAfterEndAt
	}

	return &Coupon{Entity: entity.NewEntity(), CouponCommonAttributes: couponCommonAttributes}, nil
}
