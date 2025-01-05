package orderentity

import (
	"errors"
	"time"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrDiscountMustBePositive = errors.New("discount must be positive")
	ErrStartAndEndAtRequired  = errors.New("start_at and end_at are required")
	ErrStartAtAfterEndAt      = errors.New("start_at must be before end_at")
)

type Coupon struct {
	entity.Entity
	CouponCommonAttributes
}

type CouponCommonAttributes struct {
	Discount float64
	Min      float64
	StartAt  *time.Time
	EndAt    *time.Time
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
