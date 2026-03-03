package stockusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func movementRestoreKey(m model.StockMovement) string {
	batchID := "nil"
	if m.BatchID != nil {
		batchID = m.BatchID.String()
	}

	return fmt.Sprintf("%s|%s|%s|%s", m.StockID.String(), batchID, m.GetQuantity().String(), m.GetPrice().String())
}

func getPendingOutMovementsToRestore(movements []model.StockMovement) []model.StockMovement {
	outsByKey := map[string][]model.StockMovement{}
	restoreCountByKey := map[string]int{}

	for _, m := range movements {
		switch m.Type {
		case string(stockentity.MovementTypeOut):
			key := movementRestoreKey(m)
			outsByKey[key] = append(outsByKey[key], m)
		case string(stockentity.MovementTypeRestore):
			key := movementRestoreKey(m)
			restoreCountByKey[key]++
		}
	}

	pending := make([]model.StockMovement, 0)
	for key, outs := range outsByKey {
		skip := restoreCountByKey[key]
		if skip >= len(outs) {
			continue
		}
		pending = append(pending, outs[skip:]...)
	}

	return pending
}

func reconcileStockAfterDebit(stock *stockentity.Stock, quantity decimal.Decimal, orderID uuid.UUID) {
	if orderID != uuid.Nil {
		// Pedido com possível reserva prévia.
		// Usa o saldo real de ReservedStock para saber quanto já foi pré-debitado
		// de CurrentStock na etapa de reserva.
		if stock.ReservedStock.GreaterThanOrEqual(quantity) {
			// Reserva cobriu tudo: CurrentStock já foi decrementado. Só limpar reserva.
			stock.ReservedStock = stock.ReservedStock.Sub(quantity)
			return
		}

		// Reserva parcial ou falhou silenciosamente: decrementar a parte não reservada.
		notReserved := quantity.Sub(stock.ReservedStock)
		stock.ReservedStock = decimal.Zero
		stock.CurrentStock = stock.CurrentStock.Sub(notReserved)
		return
	}

	// Saída manual (sem reserva prévia): deduzir do saldo global agora.
	stock.CurrentStock = stock.CurrentStock.Sub(quantity)
}

// DebitStockFIFO debita o estoque seguindo a estratégia FIFO (First-In, First-Out)
func (s *Service) DebitStockFIFO(ctx context.Context, stockID uuid.UUID, quantity decimal.Decimal, orderID uuid.UUID, employeeID uuid.UUID, reason string) error {
	// 1. Iniciar transação de tenant
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// 2. Bloquear Lotes com SELECT FOR UPDATE (Pessimistic Locking)
	batchesModel, err := s.stockBatchRepo.GetActiveBatchesByStockIDForUpdate(ctx, tx, stockID.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar lotes para débito: %w", err)
	}

	remainingQuantity := quantity
	movements := []*stockentity.StockMovement{}
	var movementOrderID *uuid.UUID
	if orderID != uuid.Nil {
		movementOrderID = &orderID
	}

	// 3. Processar lotes (FIFO)
	for _, bm := range batchesModel {
		if remainingQuantity.LessThanOrEqual(decimal.Zero) {
			break
		}

		batch := bm.ToDomain()
		consumeQuantity := decimal.Min(remainingQuantity, batch.CurrentQuantity)

		// Criar movimento para este lote
		movement := &stockentity.StockMovement{
			StockMovementCommonAttributes: stockentity.StockMovementCommonAttributes{
				StockID:    stockID,
				BatchID:    &batch.ID,
				Type:       stockentity.MovementTypeOut,
				Quantity:   consumeQuantity,
				Reason:     reason,
				OrderID:    movementOrderID,
				EmployeeID: employeeID,
				Price:      batch.CostPrice, // Valor de custo na saída (COGS)
			},
		}
		movement.Entity = stockentity.NewStockMovement(movement.StockID, movement.Quantity, movement.Reason, movement.EmployeeID, movement.Price).Entity

		movements = append(movements, movement)

		// Atualizar quantidade do lote
		batch.CurrentQuantity = batch.CurrentQuantity.Sub(consumeQuantity)
		remainingQuantity = remainingQuantity.Sub(consumeQuantity)

		// Salvar lote atualizado
		bm.FromDomain(batch)
		if err := s.stockBatchRepo.UpdateBatch(ctx, tx, &bm); err != nil {
			return fmt.Errorf("erro ao atualizar lote %s: %w", batch.ID, err)
		}
	}

	// 4. Se ainda sobrou quantidade, permitir estoque negativo (usando o último lote ou criando um movimento sem lote)
	if remainingQuantity.GreaterThan(decimal.Zero) {
		// Criar movimento residual (sem lote específico ou lote "fantasma")
		movement := &stockentity.StockMovement{
			StockMovementCommonAttributes: stockentity.StockMovementCommonAttributes{
				StockID:    stockID,
				Type:       stockentity.MovementTypeOut,
				Quantity:   remainingQuantity,
				Reason:     reason + " (Estoque Negativo)",
				OrderID:    movementOrderID,
				EmployeeID: employeeID,
				Price:      decimal.Zero, // Custo zero para estoque não rastreado
			},
		}
		movement.Entity = stockentity.NewStockMovement(movement.StockID, movement.Quantity, movement.Reason, movement.EmployeeID, movement.Price).Entity
		movements = append(movements, movement)
	}

	// 5. Salvar todos os movimentos
	for _, m := range movements {
		movementModel := &model.StockMovement{}
		movementModel.FromDomain(m)
		if err := s.stockMovementRepo.CreateMovement(ctx, tx, movementModel); err != nil {
			fmt.Printf("Aviso: erro ao salvar movimento de estoque: %v\n", err)
			return nil
		}
	}

	// 6. Atualizar o estoque principal (total)
	stockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, stockID.String())
	if err != nil {
		fmt.Printf("Aviso: erro ao buscar estoque principal para atualização: %v\n", err)
		return nil
	}

	stock := stockModel.ToDomain()
	reconcileStockAfterDebit(stock, quantity, orderID)

	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, tx, stockModel); err != nil {
		fmt.Printf("Aviso: erro ao atualizar estoque principal: %v\n", err)
		return nil
	}

	// 7. Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao commitar transação de débito: %w", err)
	}

	return nil
}

// RestoreStockFromOrder restaura o estoque debitado de um pedido (ex: cancelamento de pedido finalizado)
func (s *Service) RestoreStockFromOrder(ctx context.Context, orderID uuid.UUID, employeeID uuid.UUID) error {
	// 1. Buscar todos os movimentos para este pedido
	movementsModel, err := s.stockMovementRepo.GetMovementsByOrderID(ctx, orderID.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar movimentos do pedido: %w", err)
	}

	pendingOutMovements := getPendingOutMovementsToRestore(movementsModel)

	// Se não houver saídas pendentes de restauração, nada a fazer
	if len(pendingOutMovements) == 0 {
		return nil
	}

	// 2. Iniciar transação de tenant
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	for _, mm := range pendingOutMovements {
		movement := mm.ToDomain()

		// 3. Se houver lote, devolver para o lote
		if mm.BatchID != nil {
			batchModel, err := s.stockBatchRepo.GetBatchByID(ctx, tx, mm.BatchID.String())
			if err != nil {
				fmt.Printf("Aviso: erro ao buscar lote %s para restauração: %v\n", mm.BatchID.String(), err)
				continue
			}

			batch := batchModel.ToDomain()
			batch.CurrentQuantity = batch.CurrentQuantity.Add(movement.Quantity)

			batchModel.FromDomain(batch)
			if err := s.stockBatchRepo.UpdateBatch(ctx, tx, batchModel); err != nil {
				fmt.Printf("Aviso: erro ao atualizar lote %s na restauração: %v\n", batch.ID.String(), err)
				continue
			}
		}

		// 4. Criar movimento de restauração
		restoreMovement := &stockentity.StockMovement{
			StockMovementCommonAttributes: stockentity.StockMovementCommonAttributes{
				StockID:    movement.StockID,
				BatchID:    movement.BatchID,
				Type:       stockentity.MovementTypeRestore,
				Quantity:   movement.Quantity,
				Reason:     fmt.Sprintf("Restauração Pedido %s", orderID),
				OrderID:    &orderID,
				EmployeeID: employeeID,
				Price:      movement.Price,
			},
		}
		restoreMovement.Entity = stockentity.NewStockMovement(restoreMovement.StockID, restoreMovement.Quantity, restoreMovement.Reason, restoreMovement.EmployeeID, restoreMovement.Price).Entity

		restoredModel := &model.StockMovement{}
		restoredModel.FromDomain(restoreMovement)
		if err := s.stockMovementRepo.CreateMovement(ctx, tx, restoredModel); err != nil {
			fmt.Printf("Aviso: erro ao salvar movimento de restauração: %v\n", err)
			continue
		}

		// 5. Atualizar estoque principal
		stockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, movement.StockID.String())
		if err != nil {
			fmt.Printf("Aviso: erro ao buscar estoque %s para restauração: %v\n", movement.StockID.String(), err)
			continue
		}

		stock := stockModel.ToDomain()
		// Devolver para o estoque atual (já que foi debitado)
		stock.CurrentStock = stock.CurrentStock.Add(movement.Quantity)

		stockModel.FromDomain(stock)
		if err := s.stockRepo.UpdateStock(ctx, tx, stockModel); err != nil {
			fmt.Printf("Aviso: erro ao atualizar estoque principal %s na restauração: %v\n", movement.StockID.String(), err)
			continue
		}
	}

	return tx.Commit()
}
