package stockusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *Service) AddMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementCreateDTO) (*stockdto.StockMovementDTO, error) {
	// 1. Buscar estoque
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estoque %s: %w", dtoID.ID, err)
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	employee, err := s.employeeRepo.GetEmployeeByUserID(ctx, userIDUUID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar funcionário: %w", err)
	}

	// Quantidade negativa não é permitida em AddMovementStock.
	// Para saídas manuais, use RemoveMovementStock.
	if dto.Quantity.LessThanOrEqual(decimal.Zero) {
		return nil, stockentity.ErrInvalidQuantity
	}

	// Caso contrário, é uma entrada manual (cria novo lote)
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	stock := stockModel.ToDomain()

	// 2. Criar novo lote
	var variationID uuid.UUID
	if stock.ProductVariationID != nil {
		variationID = *stock.ProductVariationID
	}

	batch := &stockentity.StockBatch{
		Entity: entity.NewEntity(),
		StockBatchCommonAttributes: stockentity.StockBatchCommonAttributes{
			StockID:            stock.ID,
			ProductVariationID: variationID,
			InitialQuantity:    dto.Quantity,
			CurrentQuantity:    dto.Quantity,
			CostPrice:          dto.Price,
			ExpiresAt:          dto.ExpiresAt,
		},
	}

	batchModel := &model.StockBatch{}
	batchModel.FromDomain(batch)
	if err := s.stockBatchRepo.CreateBatch(ctx, tx, batchModel); err != nil {
		return nil, fmt.Errorf("erro ao criar lote de estoque: %w", err)
	}

	// 3. Criar movimento de entrada
	movement := &stockentity.StockMovement{
		StockMovementCommonAttributes: stockentity.StockMovementCommonAttributes{
			StockID:    stock.ID,
			BatchID:    &batch.ID,
			Type:       stockentity.MovementTypeIn,
			Quantity:   dto.Quantity,
			Reason:     dto.Reason,
			EmployeeID: employee.ID,
			Price:      dto.Price,
		},
	}
	movement.Entity = stockentity.NewStockMovement(movement.StockID, movement.Quantity, movement.Reason, movement.EmployeeID, movement.Price).Entity

	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, tx, movementModel); err != nil {
		return nil, fmt.Errorf("erro ao salvar movimento de entrada: %w", err)
	}

	// 4. Atualizar o estoque principal
	stock.CurrentStock = stock.CurrentStock.Add(dto.Quantity)
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, tx, stockModel); err != nil {
		return nil, fmt.Errorf("erro ao atualizar estoque principal: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Retornar DTO
	movementDTO := &stockdto.StockMovementDTO{}
	movementDTO.FromDomain(movement)

	return movementDTO, nil
}

// RemoveMovementStock remove estoque manualmente via FIFO
func (s *Service) RemoveMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementRemoveDTO) (*stockdto.StockMovementDTO, error) {
	// 1. Buscar estoque
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estoque %s: %w", dtoID.ID, err)
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	employee, err := s.employeeRepo.GetEmployeeByUserID(ctx, userIDUUID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar funcionário: %w", err)
	}

	// 2. Debitar via FIFO
	if err := s.DebitStockFIFO(ctx, stockModel.ID, dto.Quantity, uuid.Nil, employee.ID, dto.Reason); err != nil {
		return nil, err
	}

	// Retornar um DTO genérico ou buscar o movimento criado (para fins de retorno de API)
	return &stockdto.StockMovementDTO{}, nil
}

// AdjustMovementStock ajusta estoque para um valor específico
func (s *Service) AdjustMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementAdjustDTO) (*stockdto.StockMovementDTO, error) {
	// 1. Buscar estoque
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estoque %s: %w", dtoID.ID, err)
	}

	stock := stockModel.ToDomain()
	difference := dto.NewStock.Sub(stock.CurrentStock)

	if difference.IsZero() {
		return nil, nil
	}

	// 2. Se a diferença for negativa, é uma saída manual via FIFO
	if difference.LessThan(decimal.Zero) {
		removeDTO := &stockdto.StockMovementRemoveDTO{
			Reason:   dto.Reason + " (Ajuste)",
			Quantity: difference.Abs(),
		}
		return s.RemoveMovementStock(ctx, dtoID, removeDTO)
	}

	// 3. Se a diferença for positiva, é uma entrada manual (cria novo lote com custo zero por ser ajuste)
	addDTO := &stockdto.StockMovementCreateDTO{
		Reason:   dto.Reason + " (Ajuste)",
		Quantity: difference,
		Price:    decimal.Zero, // Custo zero para ajustes manuais de entrada se não especificado
	}
	return s.AddMovementStock(ctx, dtoID, addDTO)
}

// GetMovementsByStockID busca movimentos por estoque
func (s *Service) GetMovementsByStockID(ctx context.Context, stockID string, date *string) ([]stockdto.StockMovementDTO, error) {
	movementsModel, err := s.stockMovementRepo.GetMovementsByStockID(ctx, stockID, date)
	if err != nil {
		return nil, err
	}

	var movementsDTO []stockdto.StockMovementDTO
	for _, movementModel := range movementsModel {
		movement := movementModel.ToDomain()
		movementDTO := stockdto.StockMovementDTO{}
		movementDTO.FromDomain(movement)
		movementsDTO = append(movementsDTO, movementDTO)
	}

	return movementsDTO, nil
}
