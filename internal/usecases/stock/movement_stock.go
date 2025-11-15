package stockusecases

import (
	"context"
	"fmt"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *Service) AddMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementCreateDTO) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Adicionar estoque
	movement, err := stock.AddMovementStock(
		dto.Quantity,
		dto.Reason,
		dto.EmployeeID,
		dto.Price,
		dto.TotalPrice,
	)
	if err != nil {
		return nil, err
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
		return nil, err
	}

	// Atualizar estoque
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
		return nil, err
	}

	// Verificar alertas
	alerts := stock.CheckAlerts()
	for _, alert := range alerts {
		alertModel := &model.StockAlert{}
		alertModel.FromDomain(alert)
		if err := s.stockAlertRepo.CreateAlert(ctx, alertModel); err != nil {
			fmt.Printf("Erro ao criar alerta: %v\n", err)
		}
	}

	// Retornar DTO
	movementDTO := &stockdto.StockMovementDTO{}
	movementDTO.FromDomain(movement)

	return movementDTO, nil
}

// RemoveMovementStock remove estoque manualmente
func (s *Service) RemoveMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementRemoveDTO) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Remover estoque
	movement, err := stock.RemoveMovementStock(
		dto.Quantity,
		dto.Reason,
		dto.EmployeeID,
		dto.Price,
		dto.TotalPrice,
	)
	if err != nil {
		return nil, err
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
		return nil, err
	}

	// Atualizar estoque
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
		return nil, err
	}

	// Verificar alertas
	alerts := stock.CheckAlerts()
	for _, alert := range alerts {
		alertModel := &model.StockAlert{}
		alertModel.FromDomain(alert)
		if err := s.stockAlertRepo.CreateAlert(ctx, alertModel); err != nil {
			fmt.Printf("Erro ao criar alerta: %v\n", err)
		}
	}

	// Retornar DTO
	movementDTO := &stockdto.StockMovementDTO{}
	movementDTO.FromDomain(movement)

	return movementDTO, nil
}

// AdjustMovementStock ajusta estoque para um valor específico
func (s *Service) AdjustMovementStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockMovementAdjustDTO) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Ajustar estoque
	movement, err := stock.AdjustMovementStock(
		dto.NewStock,
		dto.Reason,
		dto.EmployeeID,
	)
	if err != nil {
		return nil, err
	}

	// Se não houve movimento (mesmo valor), retornar nil
	if movement == nil {
		return nil, nil
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
		return nil, err
	}

	// Atualizar estoque
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
		return nil, err
	}

	// Verificar alertas
	alerts := stock.CheckAlerts()
	for _, alert := range alerts {
		alertModel := &model.StockAlert{}
		alertModel.FromDomain(alert)
		if err := s.stockAlertRepo.CreateAlert(ctx, alertModel); err != nil {
			fmt.Printf("Erro ao criar alerta: %v\n", err)
		}
	}

	// Retornar DTO
	movementDTO := &stockdto.StockMovementDTO{}
	movementDTO.FromDomain(movement)

	return movementDTO, nil
}

// GetMovementsByStockID busca movimentos por estoque
func (s *Service) GetMovementsByStockID(ctx context.Context, stockID string) ([]stockdto.StockMovementDTO, error) {
	movementsModel, err := s.stockMovementRepo.GetMovementsByStockID(ctx, stockID)
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
