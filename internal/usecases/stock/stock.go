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
	itemRepo          model.ItemRepository
}

func NewStockService(
	stockRepo model.StockRepository,
	stockMovementRepo model.StockMovementRepository,
	stockAlertRepo model.StockAlertRepository,
) *Service {
	return &Service{
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		stockAlertRepo:    stockAlertRepo,
	}
}

func (s *Service) AddDependencies(productRepo model.ProductRepository, itemRepo model.ItemRepository) {
	s.itemRepo = itemRepo
	s.productRepo = productRepo
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
	dto.UpdateDomain(stock)

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
func (s *Service) GetAllStocks(ctx context.Context, page, perPage int) ([]stockdto.StockDTO, int, error) {
	stocksModel, count, err := s.stockRepo.GetAllStocks(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	var stocksDTO []stockdto.StockDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		stocksDTO = append(stocksDTO, stockDTO)
	}

	return stocksDTO, count, nil
}

// AddMovementStock adiciona estoque manualmente

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
	stockWithProduct := &stockdto.StockWithProductDTO{}

	productModel, err := s.productRepo.GetProductById(ctx, stock.ProductID.String())
	if err != nil {
		stockWithProduct.FromDomain(stock, nil)
	}

	if productModel != nil {
		product := productModel.ToDomain()
		stockWithProduct.FromDomain(stock, product)
	}

	return stockWithProduct, nil
}

// GetAllStocksWithProduct busca todos os estoques com informações dos produtos
func (s *Service) GetAllStocksWithProduct(ctx context.Context) ([]stockdto.StockWithProductDTO, error) {
	stocksModel, _, err := s.stockRepo.GetAllStocks(ctx, 1, 1000) // Fetch all (or a large page) for now as this seems to be used for reports
	if err != nil {
		return nil, err
	}

	var stocksWithProduct []stockdto.StockWithProductDTO
	for _, stockModel := range stocksModel {
		stock := stockModel.ToDomain()
		stockDTO := stockdto.StockDTO{}
		stockDTO.FromDomain(stock)

		stockWithProduct := stockdto.StockWithProductDTO{}

		// Buscar informações do produto
		productModel, err := s.productRepo.GetProductById(ctx, stock.ProductID.String())
		if err != nil {
			stockWithProduct.FromDomain(stock, nil)
		}

		if productModel != nil {
			product := productModel.ToDomain()
			stockWithProduct.FromDomain(stock, product)
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
	alertModels, err := s.stockAlertRepo.GetAllAlerts(ctx)
	if err != nil {
		return nil, err
	}

	var alertDTOs []stockdto.StockAlertDTO
	for _, alertModel := range alertModels {
		alert := alertModel.ToDomain()
		alertDTO := &stockdto.StockAlertDTO{}
		alertDTO.FromDomain(alert)
		alertDTOs = append(alertDTOs, *alertDTO)
	}

	return alertDTOs, nil
}

func (s *Service) GetAlertByID(ctx context.Context, alertID string) (*stockdto.StockAlertDTO, error) {
	alertModel, err := s.stockAlertRepo.GetAlertByID(ctx, alertID)
	if err != nil {
		return nil, err
	}

	alert := alertModel.ToDomain()
	alertDTO := &stockdto.StockAlertDTO{}
	alertDTO.FromDomain(alert)
	return alertDTO, nil
}

func (s *Service) ResolveAlert(ctx context.Context, alertID string) error {
	alert, err := s.stockAlertRepo.GetAlertByID(ctx, alertID)
	if err != nil {
		return err
	}

	if alert.IsResolved {
		return fmt.Errorf("alert already resolved")
	}

	alert.IsResolved = true
	now := time.Now()
	alert.ResolvedAt = &now

	return s.stockAlertRepo.UpdateAlert(ctx, alert)
}

func (s *Service) DeleteAlert(ctx context.Context, alertID string) error {
	return s.stockAlertRepo.DeleteAlert(ctx, alertID)
}

func (s *Service) GetStockReport(ctx context.Context, page, perPage int) (*stockdto.StockReportCompleteDTO, int, error) {
	// Buscar todos os estoques
	stockModels, count, err := s.stockRepo.GetAllStocks(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	// Buscar produtos com estoque baixo
	lowStockProducts, err := s.stockRepo.GetLowStockProducts(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Buscar produtos sem estoque
	outOfStockProducts, err := s.stockRepo.GetOutOfStockProducts(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Buscar alertas ativos
	activeAlerts, err := s.stockAlertRepo.GetActiveAlerts(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Calcular estatísticas
	totalProducts := count
	totalLowStock := len(lowStockProducts)
	totalOutOfStock := len(outOfStockProducts)
	totalActiveAlerts := len(activeAlerts)

	// Calcular valor total do estoque
	totalStockValue := decimal.Zero
	for _, stock := range stockModels {
		// Usar o custo unitário do último movimento ou um valor padrão
		// Por enquanto, vamos usar um valor estimado baseado no estoque atual
		estimatedCost := decimal.NewFromFloat(10.0) // Valor estimado por unidade
		totalStockValue = totalStockValue.Add(estimatedCost.Mul(stock.CurrentStock))
	}

	// Converter para DTOs
	var stockDTOs []stockdto.StockDTO
	for _, stockModel := range stockModels {
		stock := stockModel.ToDomain()
		stockDTO := &stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		stockDTOs = append(stockDTOs, *stockDTO)
	}

	var lowStockDTOs []stockdto.StockDTO
	for _, stockModel := range lowStockProducts {
		stock := stockModel.ToDomain()
		stockDTO := &stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		lowStockDTOs = append(lowStockDTOs, *stockDTO)
	}

	var outOfStockDTOs []stockdto.StockDTO
	for _, stockModel := range outOfStockProducts {
		stock := stockModel.ToDomain()
		stockDTO := &stockdto.StockDTO{}
		stockDTO.FromDomain(stock)
		outOfStockDTOs = append(outOfStockDTOs, *stockDTO)
	}

	var alertDTOs []stockdto.StockAlertDTO
	for _, alertModel := range activeAlerts {
		alert := alertModel.ToDomain()
		alertDTO := &stockdto.StockAlertDTO{}
		alertDTO.FromDomain(alert)
		alertDTOs = append(alertDTOs, *alertDTO)
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

	return report, count, nil
}
