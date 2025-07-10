package stockusecases

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	stockRepo         model.StockRepository
	stockMovementRepo model.StockMovementRepository
	stockAlertRepo    model.StockAlertRepository
	productRepo       model.ProductRepository
}

func NewStockService(
	stockRepo model.StockRepository,
	stockMovementRepo model.StockMovementRepository,
	stockAlertRepo model.StockAlertRepository,
	productRepo model.ProductRepository,
) *Service {
	return &Service{
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		stockAlertRepo:    stockAlertRepo,
		productRepo:       productRepo,
	}
}

// CreateStock cria um novo controle de estoque
func (s *Service) CreateStock(ctx context.Context, dto *stockdto.StockCreateDTO) (*stockdto.StockDTO, error) {
	// Verificar se o produto existe
	_, err := s.productRepo.GetProductById(ctx, dto.ProductID.String())
	if err != nil {
		return nil, fmt.Errorf("produto não encontrado: %w", err)
	}

	// Verificar se já existe estoque para este produto
	existingStock, _ := s.stockRepo.GetStockByProductID(ctx, dto.ProductID.String())
	if existingStock != nil {
		return nil, fmt.Errorf("já existe controle de estoque para este produto")
	}

	// Criar estoque
	stock := dto.ToDomain()
	stockModel := &model.Stock{}
	stockModel.FromDomain(stock)

	if err := s.stockRepo.CreateStock(ctx, stockModel); err != nil {
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
	stockDTO := &stockdto.StockDTO{}
	stockDTO.FromDomain(stock)

	return stockDTO, nil
}

// UpdateStock atualiza o estoque
func (s *Service) UpdateStock(ctx context.Context, dtoID *entitydto.IDRequest, dto *stockdto.StockUpdateDTO) error {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	stock := stockModel.ToDomain()

	// Atualizar campos
	if dto.CurrentStock != nil {
		stock.CurrentStock = *dto.CurrentStock
	}
	if dto.MinStock != nil {
		stock.MinStock = *dto.MinStock
	}
	if dto.MaxStock != nil {
		stock.MaxStock = *dto.MaxStock
	}
	if dto.Unit != nil {
		stock.Unit = *dto.Unit
	}
	if dto.IsActive != nil {
		stock.IsActive = *dto.IsActive
	}

	stockModel.FromDomain(stock)
	return s.stockRepo.UpdateStock(ctx, stockModel)
}

// GetStockByID busca estoque por ID
func (s *Service) GetStockByID(ctx context.Context, dtoID *entitydto.IDRequest) (*stockdto.StockDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()
	stockDTO := &stockdto.StockDTO{}
	stockDTO.FromDomain(stock)

	return stockDTO, nil
}

// GetStockByProductID busca estoque por produto
func (s *Service) GetStockByProductID(ctx context.Context, productID string) (*stockdto.StockDTO, error) {
	stockModel, err := s.stockRepo.GetStockByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()
	stockDTO := &stockdto.StockDTO{}
	stockDTO.FromDomain(stock)

	return stockDTO, nil
}

// GetAllStocks busca todos os estoques
func (s *Service) GetAllStocks(ctx context.Context) ([]stockdto.StockDTO, error) {
	stocksModel, err := s.stockRepo.GetAllStocks(ctx)
	if err != nil {
		return nil, err
	}

	var stocksDTO []stockdto.StockDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		stocksDTO = append(stocksDTO, stockDTO)
	}

	return stocksDTO, nil
}

// AddStock adiciona estoque manualmente
func (s *Service) AddStock(ctx context.Context, dto *stockdto.StockMovementCreateDTO) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dto.StockID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Adicionar estoque
	movement, err := stock.AddStock(
		dto.Quantity,
		dto.Reason,
		dto.EmployeeID,
		dto.UnitCost,
		dto.Notes,
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

// RemoveStock remove estoque manualmente
func (s *Service) RemoveStock(ctx context.Context, dto *stockdto.StockMovementCreateDTO) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dto.StockID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Remover estoque
	movement, err := stock.RemoveStock(
		dto.Quantity,
		dto.Reason,
		dto.EmployeeID,
		dto.UnitCost,
		dto.Notes,
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

// AdjustStock ajusta estoque para um valor específico
func (s *Service) AdjustStock(ctx context.Context, dto *stockdto.StockMovementCreateDTO, newStock decimal.Decimal) (*stockdto.StockMovementDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dto.StockID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()

	// Ajustar estoque
	movement, err := stock.AdjustStock(
		newStock,
		dto.Reason,
		dto.EmployeeID,
		dto.UnitCost,
		dto.Notes,
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

// GetLowStockProducts busca produtos com estoque baixo
func (s *Service) GetLowStockProducts(ctx context.Context) ([]stockdto.StockDTO, error) {
	stocksModel, err := s.stockRepo.GetLowStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	var stocksDTO []stockdto.StockDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		stocksDTO = append(stocksDTO, stockDTO)
	}

	return stocksDTO, nil
}

// GetOutOfStockProducts busca produtos sem estoque
func (s *Service) GetOutOfStockProducts(ctx context.Context) ([]stockdto.StockDTO, error) {
	stocksModel, err := s.stockRepo.GetOutOfStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	var stocksDTO []stockdto.StockDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		stocksDTO = append(stocksDTO, stockDTO)
	}

	return stocksDTO, nil
}

// GetStockWithProduct busca estoque com informações do produto
func (s *Service) GetStockWithProduct(ctx context.Context, dtoID *entitydto.IDRequest) (*stockdto.StockWithProductDTO, error) {
	stockModel, err := s.stockRepo.GetStockByID(ctx, dtoID.ID.String())
	if err != nil {
		return nil, err
	}

	stock := stockModel.ToDomain()
	stockDTO := &stockdto.StockDTO{}
	stockDTO.FromDomain(stock)

	// Buscar informações do produto
	stockWithProduct := &stockdto.StockWithProductDTO{
		StockDTO: *stockDTO,
	}

	if s.productRepo != nil {
		product, err := s.productRepo.GetProductById(ctx, stock.ProductID.String())
		if err == nil && product != nil {
			stockWithProduct.ProductName = product.Name
			stockWithProduct.ProductCode = product.Code
		}
	}

	return stockWithProduct, nil
}

// GetAllStocksWithProduct busca todos os estoques com informações dos produtos
func (s *Service) GetAllStocksWithProduct(ctx context.Context) ([]stockdto.StockWithProductDTO, error) {
	stocksModel, err := s.stockRepo.GetAllStocks(ctx)
	if err != nil {
		return nil, err
	}

	var stocksWithProduct []stockdto.StockWithProductDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)

		stockWithProduct := stockdto.StockWithProductDTO{
			StockDTO: stockDTO,
		}

		// Buscar informações do produto
		if s.productRepo != nil {
			product, err := s.productRepo.GetProductById(ctx, stock.ProductID.String())
			if err == nil && product != nil {
				stockWithProduct.ProductName = product.Name
				stockWithProduct.ProductCode = product.Code
			}
		}

		stocksWithProduct = append(stocksWithProduct, stockWithProduct)
	}

	return stocksWithProduct, nil
}

// GetActiveAlerts busca alertas ativos de estoque
func (s *Service) GetActiveAlerts(ctx context.Context) ([]stockdto.StockAlertDTO, error) {
	alertsModel, err := s.stockAlertRepo.GetActiveAlerts(ctx)
	if err != nil {
		return nil, err
	}

	var alertsDTO []stockdto.StockAlertDTO
	for _, alertModel := range alertsModel {
		alert := alertModel.ToDomain()
		alertDTO := stockdto.StockAlertDTO{}
		alertDTO.FromDomain(alert)
		alertsDTO = append(alertsDTO, alertDTO)
	}

	return alertsDTO, nil
}

func (s *Service) GetAllAlerts(ctx context.Context) ([]stockdto.StockAlertDTO, error) {
	alerts, err := s.stockAlertRepo.GetAllAlerts(ctx)
	if err != nil {
		return nil, err
	}

	var alertDTOs []stockdto.StockAlertDTO
	for _, alert := range alerts {
		alertDTOs = append(alertDTOs, *alert.ToDTO())
	}

	return alertDTOs, nil
}

func (s *Service) GetAlertByID(ctx context.Context, alertID string) (*stockdto.StockAlertDTO, error) {
	alert, err := s.stockAlertRepo.GetAlertByID(ctx, alertID)
	if err != nil {
		return nil, err
	}

	return alert.ToDTO(), nil
}

func (s *Service) ResolveAlert(ctx context.Context, alertID string) error {
	alert, err := s.stockAlertRepo.GetAlertByID(ctx, alertID)
	if err != nil {
		return err
	}

	alert.IsResolved = true
	now := time.Now()
	alert.ResolvedAt = &now

	return s.stockAlertRepo.UpdateAlert(ctx, alert)
}

func (s *Service) DeleteAlert(ctx context.Context, alertID string) error {
	return s.stockAlertRepo.DeleteAlert(ctx, alertID)
}

func (s *Service) GetStockReport(ctx context.Context) (*stockdto.StockReportCompleteDTO, error) {
	// Buscar todos os estoques
	stocks, err := s.stockRepo.GetAllStocks(ctx)
	if err != nil {
		return nil, err
	}

	// Buscar produtos com estoque baixo
	lowStockProducts, err := s.stockRepo.GetLowStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Buscar produtos sem estoque
	outOfStockProducts, err := s.stockRepo.GetOutOfStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Buscar alertas ativos
	activeAlerts, err := s.stockAlertRepo.GetActiveAlerts(ctx)
	if err != nil {
		return nil, err
	}

	// Calcular estatísticas
	totalProducts := len(stocks)
	totalLowStock := len(lowStockProducts)
	totalOutOfStock := len(outOfStockProducts)
	totalActiveAlerts := len(activeAlerts)

	// Calcular valor total do estoque
	totalStockValue := decimal.Zero
	for _, stock := range stocks {
		// Usar o custo unitário do último movimento ou um valor padrão
		// Por enquanto, vamos usar um valor estimado baseado no estoque atual
		estimatedCost := decimal.NewFromFloat(10.0) // Valor estimado por unidade
		totalStockValue = totalStockValue.Add(estimatedCost.Mul(stock.CurrentStock))
	}

	// Converter para DTOs
	var stockDTOs []stockdto.StockDTO
	for _, stock := range stocks {
		stockDTO := stock.ToDTO()
		if stockDTO != nil {
			stockDTOs = append(stockDTOs, *stockDTO)
		}
	}

	var lowStockDTOs []stockdto.StockDTO
	for _, stock := range lowStockProducts {
		stockDTO := stock.ToDTO()
		if stockDTO != nil {
			lowStockDTOs = append(lowStockDTOs, *stockDTO)
		}
	}

	var outOfStockDTOs []stockdto.StockDTO
	for _, stock := range outOfStockProducts {
		stockDTO := stock.ToDTO()
		if stockDTO != nil {
			outOfStockDTOs = append(outOfStockDTOs, *stockDTO)
		}
	}

	var alertDTOs []stockdto.StockAlertDTO
	for _, alert := range activeAlerts {
		alertDTO := alert.ToDTO()
		if alertDTO != nil {
			alertDTOs = append(alertDTOs, *alertDTO)
		}
	}

	report := &stockdto.StockReportCompleteDTO{
		Summary: stockdto.StockReportSummaryDTO{
			TotalProducts:     totalProducts,
			TotalLowStock:     totalLowStock,
			TotalOutOfStock:   totalOutOfStock,
			TotalActiveAlerts: totalActiveAlerts,
			TotalStockValue:   totalStockValue,
		},
		AllStocks:          stockDTOs,
		LowStockProducts:   lowStockDTOs,
		OutOfStockProducts: outOfStockDTOs,
		ActiveAlerts:       alertDTOs,
		GeneratedAt:        time.Now(),
	}

	return report, nil
}
