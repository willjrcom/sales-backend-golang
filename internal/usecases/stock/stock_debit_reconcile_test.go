package stockusecases

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestReconcileStockAfterDebit_OrderWithFullReserve(t *testing.T) {
	stock := newStock(10, 0, 100)
	stock.ReservedStock = decimal.NewFromInt(5)

	reconcileStockAfterDebit(stock, decimal.NewFromInt(5), uuid.New())

	if stock.CurrentStock.String() != "10" {
		t.Fatalf("expected current_stock unchanged when fully reserved, got %s", stock.CurrentStock.String())
	}
	if stock.ReservedStock.String() != "0" {
		t.Fatalf("expected reserved_stock to be cleared, got %s", stock.ReservedStock.String())
	}
}

func TestReconcileStockAfterDebit_OrderWithoutReserveFallbackDebitsCurrentStock(t *testing.T) {
	stock := newStock(10, 0, 100)
	stock.ReservedStock = decimal.Zero

	reconcileStockAfterDebit(stock, decimal.NewFromInt(5), uuid.New())

	if stock.CurrentStock.String() != "5" {
		t.Fatalf("expected current_stock to be debited when reserve is zero, got %s", stock.CurrentStock.String())
	}
	if stock.ReservedStock.String() != "0" {
		t.Fatalf("expected reserved_stock to remain zero, got %s", stock.ReservedStock.String())
	}
}

func TestReconcileStockAfterDebit_OrderWithPartialReserveDebitsOnlyUnreserved(t *testing.T) {
	stock := newStock(10, 0, 100)
	stock.ReservedStock = decimal.NewFromInt(2)

	reconcileStockAfterDebit(stock, decimal.NewFromInt(5), uuid.New())

	if stock.CurrentStock.String() != "7" {
		t.Fatalf("expected current_stock to debit only unreserved quantity, got %s", stock.CurrentStock.String())
	}
	if stock.ReservedStock.String() != "0" {
		t.Fatalf("expected reserved_stock to be cleared, got %s", stock.ReservedStock.String())
	}
}
