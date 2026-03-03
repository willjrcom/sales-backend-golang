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
				OrderID:    &orderID,
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
				OrderID:    &orderID,
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
			return fmt.Errorf("erro ao salvar movimento de estoque: %w", err)
		}
	}

	// 6. Atualizar o estoque principal (total)
	stockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, stockID.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar estoque principal para atualização: %w", err)
	}

	stock := stockModel.ToDomain()
	if orderID != uuid.Nil {
		// Pedido com possível reserva prévia.
		// Usa o saldo real de ReservedStock para saber quanto já foi pré-debitado
		// de CurrentStock na etapa de reserva.
		if stock.ReservedStock.GreaterThanOrEqual(quantity) {
			// Reserva cobriu tudo: CurrentStock já foi decrementado. Só limpar reserva.
			stock.ReservedStock = stock.ReservedStock.Sub(quantity)
		} else {
			// Reserva parcial ou falhou silenciosamente: decrementar a parte não reservada.
			notReserved := quantity.Sub(stock.ReservedStock)
			stock.ReservedStock = decimal.Zero
			stock.CurrentStock = stock.CurrentStock.Sub(notReserved)
		}
	} else {
		// Saída manual (sem reserva prévia): deduzir do saldo global agora.
		stock.CurrentStock = stock.CurrentStock.Sub(quantity)
	}

	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, tx, stockModel); err != nil {
		return fmt.Errorf("erro ao atualizar estoque principal: %w", err)
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

	// Se não houver movimentos, nada a fazer
	if len(movementsModel) == 0 {
		return nil
	}

	// 2. Iniciar transação de tenant
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	for _, mm := range movementsModel {
		// Apenas restaurar o que foi efetivamente retirado (out)
		if mm.Type != string(stockentity.MovementTypeOut) {
			continue
		}

		movement := mm.ToDomain()

		// 3. Se houver lote, devolver para o lote
		if mm.BatchID != nil {
			batchModel, err := s.stockBatchRepo.GetBatchByID(ctx, tx, mm.BatchID.String())
			if err != nil {
				return fmt.Errorf("erro ao buscar lote %s para restauração: %w", mm.BatchID.String(), err)
			}

			batch := batchModel.ToDomain()
			batch.CurrentQuantity = batch.CurrentQuantity.Add(movement.Quantity)

			batchModel.FromDomain(batch)
			if err := s.stockBatchRepo.UpdateBatch(ctx, tx, batchModel); err != nil {
				return fmt.Errorf("erro ao atualizar lote %s na restauração: %w", batch.ID.String(), err)
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
			return fmt.Errorf("erro ao salvar movimento de restauração: %w", err)
		}

		// 5. Atualizar estoque principal
		stockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, movement.StockID.String())
		if err != nil {
			return fmt.Errorf("erro ao buscar estoque %s para restauração: %w", movement.StockID.String(), err)
		}

		stock := stockModel.ToDomain()
		// Devolver para o estoque atual (já que foi debitado)
		stock.CurrentStock = stock.CurrentStock.Add(movement.Quantity)

		stockModel.FromDomain(stock)
		if err := s.stockRepo.UpdateStock(ctx, tx, stockModel); err != nil {
			return fmt.Errorf("erro ao atualizar estoque principal %s na restauração: %w", movement.StockID.String(), err)
		}
	}

	return tx.Commit()
}
