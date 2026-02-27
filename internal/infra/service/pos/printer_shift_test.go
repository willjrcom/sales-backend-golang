package pos

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

func Test_FormatShift(t *testing.T) {
	shift := shiftentity.NewShift(decimal.NewFromInt(10))
	shift.Orders = []orderentity.Order{
		{
			OrderCommonAttributes: orderentity.OrderCommonAttributes{
				OrderNumber: 1,
				OrderDetail: orderentity.OrderDetail{
					SubTotal: decimal.NewFromInt(10),
				},
			},
		},
		{
			OrderCommonAttributes: orderentity.OrderCommonAttributes{
				OrderNumber: 2,
				OrderDetail: orderentity.OrderDetail{
					SubTotal: decimal.NewFromInt(131),
				},
			},
		},
		{
			OrderCommonAttributes: orderentity.OrderCommonAttributes{
				OrderNumber: 3,
				OrderDetail: orderentity.OrderDetail{
					SubTotal: decimal.NewFromInt(110),
				},
			},
		},
		{
			OrderCommonAttributes: orderentity.OrderCommonAttributes{
				OrderNumber: 4,
				OrderDetail: orderentity.OrderDetail{
					SubTotal: decimal.NewFromInt(10),
				},
			},
		},
	}

	out, err := FormatShift(shift)
	assert.NoError(t, err)
	if err := os.WriteFile("printer_shift.txt", out, 0644); err != nil {
		t.Fatalf("failed to write printer buffer to file: %v", err)
	}
}
