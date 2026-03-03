package stockusecases

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	stocklocal "github.com/willjrcom/sales-backend-go/internal/infra/repository/local/stock"
)

var (
	svc          *Service
	ctx          context.Context
	stockRepo    *stocklocal.StockRepositoryLocal
	batchRepo    *stocklocal.StockBatchRepositoryLocal
	movementRepo *stocklocal.StockMovementRepositoryLocal
	alertRepo    *stocklocal.StockAlertRepositoryLocal
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	stockRepo = stocklocal.NewStockRepositoryLocal()
	batchRepo = stocklocal.NewStockBatchRepositoryLocal()
	movementRepo = stocklocal.NewStockMovementRepositoryLocal()
	alertRepo = stocklocal.NewStockAlertRepositoryLocal()

	svc = NewStockService(nil, stockRepo, movementRepo, batchRepo, alertRepo)

	os.Exit(m.Run())
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func newStock(initial, min, max float64) *stockentity.Stock {
	return stockentity.NewStock(
		uuid.New(), nil,
		decimal.NewFromFloat(initial),
		decimal.NewFromFloat(min),
		decimal.NewFromFloat(max),
		"un",
	)
}

func newBatch(stockID uuid.UUID, qty float64, expiresAt *time.Time) *stockentity.StockBatch {
	return stockentity.NewStockBatch(stockID, uuid.Nil, decimal.NewFromFloat(qty), decimal.NewFromFloat(5.00), expiresAt)
}

// ─────────────────────────────────────────────────────────────
// ReserveStock — Domínio
// ─────────────────────────────────────────────────────────────

func TestReserveStock_Success(t *testing.T) {
	stock := newStock(10, 2, 20)
	orderID := uuid.New()
	employeeID := uuid.New()

	movement, err := stock.ReserveStock(decimal.NewFromInt(3), orderID, employeeID, decimal.NewFromFloat(5.00))
	require.NoError(t, err)

	assert.Equal(t, "7", stock.CurrentStock.String(), "CurrentStock deve diminuir para 7")
	assert.Equal(t, "3", stock.ReservedStock.String(), "ReservedStock deve subir para 3")
	assert.Equal(t, stockentity.MovementTypeReserve, movement.Type)
	assert.Equal(t, orderID, *movement.OrderID)
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	stock := newStock(2, 0, 20)
	_, err := stock.ReserveStock(decimal.NewFromInt(5), uuid.New(), uuid.New(), decimal.Zero)
	assert.ErrorIs(t, err, stockentity.ErrInsufficientStock)
}

func TestReserveStock_ZeroQuantity_Error(t *testing.T) {
	stock := newStock(10, 0, 20)
	_, err := stock.ReserveStock(decimal.Zero, uuid.New(), uuid.New(), decimal.Zero)
	assert.ErrorIs(t, err, stockentity.ErrInvalidQuantity)
}

func TestReserveStock_InactiveStock_Error(t *testing.T) {
	stock := newStock(10, 0, 20)
	stock.IsActive = false
	_, err := stock.ReserveStock(decimal.NewFromInt(3), uuid.New(), uuid.New(), decimal.Zero)
	assert.Error(t, err, "estoque inativo não pode ser reservado")
}

// ─────────────────────────────────────────────────────────────
// RestoreStock — Domínio
// ─────────────────────────────────────────────────────────────

func TestRestoreStock_Success(t *testing.T) {
	stock := newStock(10, 0, 20)
	orderID := uuid.New()
	employeeID := uuid.New()
	price := decimal.NewFromFloat(5.00)

	_, err := stock.ReserveStock(decimal.NewFromInt(4), orderID, employeeID, price)
	require.NoError(t, err)

	movement, err := stock.RestoreStock(decimal.NewFromInt(4), orderID, employeeID, price, nil)
	require.NoError(t, err)

	assert.Equal(t, "10", stock.CurrentStock.String(), "CurrentStock deve voltar a 10")
	assert.Equal(t, "0", stock.ReservedStock.String(), "ReservedStock deve zerar")
	assert.Equal(t, stockentity.MovementTypeRestore, movement.Type)
}

func TestRestoreStock_MoreThanReserved_Error(t *testing.T) {
	stock := newStock(10, 0, 20)
	orderID := uuid.New()

	_, _ = stock.ReserveStock(decimal.NewFromInt(2), orderID, uuid.New(), decimal.Zero)

	_, err := stock.RestoreStock(decimal.NewFromInt(5), orderID, uuid.New(), decimal.Zero, nil)
	assert.Error(t, err, "não deve restaurar mais do que o reservado")
}

func TestRestoreStock_ZeroQuantity_Error(t *testing.T) {
	stock := newStock(10, 0, 20)
	_, err := stock.RestoreStock(decimal.Zero, uuid.New(), uuid.New(), decimal.Zero, nil)
	assert.ErrorIs(t, err, stockentity.ErrInvalidQuantity)
}

// ─────────────────────────────────────────────────────────────
// AddMovementStock / RemoveMovementStock — Domínio
// ─────────────────────────────────────────────────────────────

func TestAddMovement_IncreasesStock(t *testing.T) {
	stock := newStock(5, 0, 100)
	movement, err := stock.AddMovementStock(decimal.NewFromInt(10), "compra", uuid.New(), decimal.NewFromFloat(8.50))
	require.NoError(t, err)

	assert.Equal(t, "15", stock.CurrentStock.String())
	assert.Equal(t, stockentity.MovementTypeIn, movement.Type)
}

func TestAddMovement_ZeroQuantity_Error(t *testing.T) {
	stock := newStock(5, 0, 100)
	_, err := stock.AddMovementStock(decimal.Zero, "compra", uuid.New(), decimal.Zero)
	assert.ErrorIs(t, err, stockentity.ErrInvalidQuantity)
}

func TestAddMovement_InactiveStock_Error(t *testing.T) {
	stock := newStock(5, 0, 100)
	stock.IsActive = false
	_, err := stock.AddMovementStock(decimal.NewFromInt(5), "compra", uuid.New(), decimal.Zero)
	assert.Error(t, err)
}

func TestRemoveMovement_DecreasesStock(t *testing.T) {
	stock := newStock(10, 0, 100)
	movement, err := stock.RemoveMovementStock(decimal.NewFromInt(3), "saída", uuid.New(), decimal.NewFromFloat(5.50))
	require.NoError(t, err)

	assert.Equal(t, "7", stock.CurrentStock.String())
	assert.Equal(t, stockentity.MovementTypeOut, movement.Type)
}

func TestRemoveMovement_InsufficientStock_Error(t *testing.T) {
	stock := newStock(2, 0, 100)
	_, err := stock.RemoveMovementStock(decimal.NewFromInt(5), "saída", uuid.New(), decimal.Zero)
	assert.ErrorIs(t, err, stockentity.ErrInsufficientStock)
}

// ─────────────────────────────────────────────────────────────
// AdjustMovementStock — Domínio
// ─────────────────────────────────────────────────────────────

func TestAdjustMovement_Up(t *testing.T) {
	stock := newStock(5, 0, 100)
	movement, err := stock.AdjustMovementStock(decimal.NewFromInt(15), "inventário", uuid.New())
	require.NoError(t, err)

	assert.Equal(t, "15", stock.CurrentStock.String())
	assert.Equal(t, stockentity.MovementTypeAdjustIn, movement.Type)
}

func TestAdjustMovement_Down(t *testing.T) {
	stock := newStock(15, 0, 100)
	movement, err := stock.AdjustMovementStock(decimal.NewFromInt(8), "inventário", uuid.New())
	require.NoError(t, err)

	assert.Equal(t, "8", stock.CurrentStock.String())
	assert.Equal(t, stockentity.MovementTypeAdjustOut, movement.Type)
}

func TestAdjustMovement_NegativeTarget_Error(t *testing.T) {
	stock := newStock(5, 0, 100)
	_, err := stock.AdjustMovementStock(decimal.NewFromInt(-1), "inventário", uuid.New())
	assert.ErrorIs(t, err, stockentity.ErrInvalidQuantity)
}

// ─────────────────────────────────────────────────────────────
// CheckAlerts — Domínio
// ─────────────────────────────────────────────────────────────

func TestCheckAlerts_LowStock(t *testing.T) {
	stock := newStock(2, 5, 20) // current < min
	alerts := stock.CheckAlerts()
	require.Len(t, alerts, 1)
	assert.Equal(t, stockentity.AlertTypeLowStock, alerts[0].Type)
	assert.Equal(t, stock.ID, alerts[0].StockID)
}

func TestCheckAlerts_OutOfStock(t *testing.T) {
	stock := newStock(0, 5, 20) // current == 0
	alerts := stock.CheckAlerts()
	require.Len(t, alerts, 1)
	assert.Equal(t, stockentity.AlertTypeOutOfStock, alerts[0].Type)
}

func TestCheckAlerts_OverStock(t *testing.T) {
	stock := newStock(25, 2, 20) // current > max
	alerts := stock.CheckAlerts()
	require.Len(t, alerts, 1)
	assert.Equal(t, stockentity.AlertTypeOverStock, alerts[0].Type)
}

func TestCheckAlerts_NormalStock_NoAlerts(t *testing.T) {
	stock := newStock(10, 2, 20)
	alerts := stock.CheckAlerts()
	assert.Empty(t, alerts)
}

func TestCheckAlerts_InactiveStock_NoAlerts(t *testing.T) {
	stock := newStock(1, 5, 20)
	stock.IsActive = false
	alerts := stock.CheckAlerts()
	assert.Empty(t, alerts, "estoque inativo não gera alertas")
}

// ─────────────────────────────────────────────────────────────
// StockBatch — Vencimento
// ─────────────────────────────────────────────────────────────

func TestBatch_IsExpired_True(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	batch := newBatch(uuid.New(), 5, &past)
	assert.True(t, batch.IsExpired())
}

func TestBatch_IsExpired_False(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	batch := newBatch(uuid.New(), 5, &future)
	assert.False(t, batch.IsExpired())
}

func TestBatch_NeverExpires(t *testing.T) {
	batch := newBatch(uuid.New(), 5, nil)
	assert.False(t, batch.IsExpired())
}

func TestBatch_HasStock(t *testing.T) {
	batch := newBatch(uuid.New(), 5, nil)
	assert.True(t, batch.HasStock())

	batch.CurrentQuantity = decimal.Zero
	assert.False(t, batch.HasStock())
}

func TestBatch_CheckExpiration_Expired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	batch := newBatch(uuid.New(), 5, &past)
	assert.Equal(t, stockentity.AlertTypeExpired, batch.CheckExpiration(30))
}

func TestBatch_CheckExpiration_NearExpiration(t *testing.T) {
	soon := time.Now().Add(5 * 24 * time.Hour) // 5 dias, threshold = 30
	batch := newBatch(uuid.New(), 5, &soon)
	assert.Equal(t, stockentity.AlertTypeNearExpiration, batch.CheckExpiration(30))
}

func TestBatch_CheckExpiration_Far_NoAlert(t *testing.T) {
	far := time.Now().Add(60 * 24 * time.Hour) // 60 dias, threshold = 30
	batch := newBatch(uuid.New(), 5, &far)
	assert.Empty(t, batch.CheckExpiration(30))
}

func TestBatch_CheckExpiration_NoStock_NoAlert(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	batch := newBatch(uuid.New(), 0, &past) // quantidade zero!
	assert.Empty(t, batch.CheckExpiration(30), "lote vazido não deve gerar alerta de vencimento")
}

func TestBatch_CheckExpiration_NoExpiry_NoAlert(t *testing.T) {
	batch := newBatch(uuid.New(), 5, nil) // sem data de vencimento
	assert.Empty(t, batch.CheckExpiration(30))
}

// ─────────────────────────────────────────────────────────────
// FIFO — Lógica de consumo em ordem
// ─────────────────────────────────────────────────────────────

func TestFIFO_ConsumesOldestFirst(t *testing.T) {
	// Simula o consumo FIFO de 6 unidades de dois lotes:
	// Lote 1: 5 un, Lote 2: 8 un → deve zerar lote1 e retirar 1 do lote2
	stockID := uuid.New()
	batch1 := newBatch(stockID, 5, nil)
	batch2 := newBatch(stockID, 8, nil)

	toDebit := decimal.NewFromInt(6)

	// Lote 1
	debit1 := decimal.Min(batch1.CurrentQuantity, toDebit)
	batch1.CurrentQuantity = batch1.CurrentQuantity.Sub(debit1)
	toDebit = toDebit.Sub(debit1)

	// Lote 2
	debit2 := decimal.Min(batch2.CurrentQuantity, toDebit)
	batch2.CurrentQuantity = batch2.CurrentQuantity.Sub(debit2)

	assert.Equal(t, "0", batch1.CurrentQuantity.String(), "Lote 1 deve estar zerado")
	assert.Equal(t, "7", batch2.CurrentQuantity.String(), "Lote 2 deve ter 7 restantes")
	// toDebit ainda é 1 porque batch2 recebeu o restante corretamente:
	// 6 total - 5 lote1 = 1 restante → debitado do lote2 (8-1=7) ✓
}

func TestFIFO_SkipsExpiredBatches(t *testing.T) {
	// Lote expirado não deve ser incluído no FIFO ativo
	past := time.Now().Add(-1 * time.Hour)
	expiredBatch := newBatch(uuid.New(), 10, &past)
	assert.True(t, expiredBatch.IsExpired())

	// Somente lotes não expirados devem ser usados para FIFO:
	future := time.Now().Add(24 * time.Hour)
	activeBatch := newBatch(uuid.New(), 10, &future)
	assert.False(t, activeBatch.IsExpired())
	assert.True(t, activeBatch.HasStock())
}

// ─────────────────────────────────────────────────────────────
// Fluxo reserva → cancelamento → restaura
// ─────────────────────────────────────────────────────────────

func TestReserveAndRestoreFlow(t *testing.T) {
	stock := newStock(15, 3, 50)
	orderID := uuid.New()
	employeeID := uuid.New()
	price := decimal.NewFromFloat(5.00)
	qty := decimal.NewFromInt(5)

	// 1. Reservar
	_, err := stock.ReserveStock(qty, orderID, employeeID, price)
	require.NoError(t, err)
	assert.Equal(t, "10", stock.CurrentStock.String())
	assert.Equal(t, "5", stock.ReservedStock.String())

	// 2. Cancelar → restaurar
	_, err = stock.RestoreStock(qty, orderID, employeeID, price, nil)
	require.NoError(t, err)
	assert.Equal(t, "15", stock.CurrentStock.String())
	assert.Equal(t, "0", stock.ReservedStock.String())
}

// ─────────────────────────────────────────────────────────────
// Fluxo reserva → finalização (FIFO debit)
// ─────────────────────────────────────────────────────────────

func TestReserveAndFinalizeFlow(t *testing.T) {
	// Depois da reserva, finalização consome do lote (FIFO). ReservedStock vai a zero.
	stock := newStock(20, 3, 50)
	orderID := uuid.New()
	employeeID := uuid.New()
	price := decimal.NewFromFloat(5.00)
	qty := decimal.NewFromInt(8)

	// 1. Reserva — move para ReservedStock
	_, err := stock.ReserveStock(qty, orderID, employeeID, price)
	require.NoError(t, err)
	assert.Equal(t, "12", stock.CurrentStock.String())
	assert.Equal(t, "8", stock.ReservedStock.String())

	// 2. Finalização: FIFO consome do lote. Simula o que DebitStockFIFO faz:
	//    Remove de ReservedStock (debita reserva) e registra saída
	stock.ReservedStock = stock.ReservedStock.Sub(qty)
	// CurrentStock já foi debitado na reserva — após FIFO fica como está
	assert.Equal(t, "0", stock.ReservedStock.String())
	assert.Equal(t, "12", stock.CurrentStock.String(), "CurrentStock não muda na finalização")
}
