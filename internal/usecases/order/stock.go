package orderusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *ItemService) reserveStockFromItemWithTx(ctx context.Context, tx *bun.Tx, item *orderentity.Item, orderID uuid.UUID, attendantID uuid.UUID) error {
	// Buscar estoque do produto/variação
	stockModel, err := s.getStockToMove(ctx, item)
	if stockModel == nil || err != nil {
		fmt.Printf("stock not exists for item: %s", item.Name)
		return nil
	}

	lockedStockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, stockModel.ID.String())
	if err != nil {
		return fmt.Errorf("erro ao bloquear estoque para reserva: %w", err)
	}

	stock := lockedStockModel.ToDomain()

	movement, err := stock.ReserveStock(
		decimal.NewFromFloat(item.Quantity),
		orderID,
		attendantID,
		item.SubTotal,
	)
	if err != nil {
		fmt.Printf("Aviso: erro ao reservar estoque para produto %s: %v\n", item.Name, err)
		return err
	}

	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, tx, movementModel); err != nil {
		fmt.Printf("Aviso: erro ao salvar movimento de estoque: %v\n", err)
		return err
	}

	// Atualizar estoque (CurrentStock diminuiu, ReservedStock aumentou)
	lockedStockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, tx, lockedStockModel); err != nil {
		fmt.Printf("Aviso: erro ao atualizar estoque: %v\n", err)
		return err
	}

	fmt.Printf("Estoque disponível: %f • Restará após adicionar: %f para produto %s\n", stock.CurrentStock.Add(decimal.NewFromFloat(item.Quantity)).InexactFloat64(), stock.CurrentStock.InexactFloat64(), item.Name)

	return nil
}

func (s *ItemService) RestoreStockFromItem(ctx context.Context, item *orderentity.Item, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if err := s.restoreStockFromItemWithTx(ctx, tx, item, groupItem, attendantID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ItemService) restoreStockFromItemWithTx(ctx context.Context, tx *bun.Tx, item *orderentity.Item, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	// Buscar estoque do produto/variação
	stockModel, err := s.getStockToMove(ctx, item)
	if stockModel == nil || err != nil {
		fmt.Printf("stock not exists for item: %s", item.Name)
		return nil
	}

	lockedStockModel, err := s.stockRepo.GetStockByIDForUpdate(ctx, tx, stockModel.ID.String())
	if err != nil {
		return fmt.Errorf("erro ao bloquear estoque para restauração: %w", err)
	}

	stock := lockedStockModel.ToDomain()

	movement, err := stock.RestoreStock(
		decimal.NewFromFloat(item.Quantity),
		groupItem.OrderID,
		attendantID,
		item.SubTotal,
		nil, // BatchID not available at reservation restoration
	)
	if err != nil {
		fmt.Printf("Aviso: erro ao restaurar estoque para produto %s: %v\n", item.Name, err)
		return err
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, tx, movementModel); err != nil {
		fmt.Printf("Aviso: erro ao salvar movimento de estoque: %v\n", err)
		return err
	}

	// Atualizar estoque
	lockedStockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, tx, lockedStockModel); err != nil {
		fmt.Printf("Aviso: erro ao atualizar estoque: %v\n", err)
		return err
	}

	return nil
}
