package stockusecases

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
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

func TestGetPendingOutMovementsToRestore_IgnoresAlreadyRestoredOuts(t *testing.T) {
	stockID := uuid.New()
	batchID := uuid.New()
	qty := decimal.NewFromInt(2)
	price := decimal.NewFromInt(10)

	out := model.StockMovement{}
	out.StockID = stockID
	out.BatchID = &batchID
	out.Type = string(stockentity.MovementTypeOut)
	out.Quantity = &qty
	out.Price = &price

	restore := model.StockMovement{}
	restore.StockID = stockID
	restore.BatchID = &batchID
	restore.Type = string(stockentity.MovementTypeRestore)
	restore.Quantity = &qty
	restore.Price = &price

	pending := getPendingOutMovementsToRestore([]model.StockMovement{out, restore})
	if len(pending) != 0 {
		t.Fatalf("expected no pending outs when already restored, got %d", len(pending))
	}
}

func TestGetPendingOutMovementsToRestore_OnlyReturnsNonRestoredOuts(t *testing.T) {
	stockID := uuid.New()
	qty := decimal.NewFromInt(1)
	price := decimal.NewFromInt(5)

	outA := model.StockMovement{}
	outA.StockID = stockID
	outA.Type = string(stockentity.MovementTypeOut)
	outA.Quantity = &qty
	outA.Price = &price

	outB := model.StockMovement{}
	outB.StockID = stockID
	outB.Type = string(stockentity.MovementTypeOut)
	outB.Quantity = &qty
	outB.Price = &price

	restore := model.StockMovement{}
	restore.StockID = stockID
	restore.Type = string(stockentity.MovementTypeRestore)
	restore.Quantity = &qty
	restore.Price = &price

	pending := getPendingOutMovementsToRestore([]model.StockMovement{outA, outB, restore})
	if len(pending) != 1 {
		t.Fatalf("expected one pending out after matching one restore, got %d", len(pending))
	}
}
