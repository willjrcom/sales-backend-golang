package stockusecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	db                *bun.DB
	stockRepo         model.StockRepository
	stockMovementRepo model.StockMovementRepository
	stockBatchRepo    model.StockBatchRepository
	stockAlertRepo    model.StockAlertRepository
	productRepo       model.ProductRepository
	itemRepo          model.ItemRepository
	employeeRepo      model.EmployeeRepository
	orderRepo         model.OrderRepository
}

func NewStockService(
	db *bun.DB,
	stockRepo model.StockRepository,
	stockMovementRepo model.StockMovementRepository,
	stockBatchRepo model.StockBatchRepository,
	stockAlertRepo model.StockAlertRepository,
) *Service {
	return &Service{
		db:                db,
		stockRepo:         stockRepo,
		stockMovementRepo: stockMovementRepo,
		stockBatchRepo:    stockBatchRepo,
		stockAlertRepo:    stockAlertRepo,
	}
}

func (s *Service) AddDependencies(productRepo model.ProductRepository, itemRepo model.ItemRepository, employeeRepo model.EmployeeRepository, orderRepo model.OrderRepository) {
	s.itemRepo = itemRepo
	s.productRepo = productRepo
	s.employeeRepo = employeeRepo
	s.orderRepo = orderRepo
}

// createAlertsIfNotDuplicate persiste cada alerta somente se não houver outro alerta
// ativo do mesmo tipo para o mesmo estoque (evita duplicatas).
func (s *Service) createAlertsIfNotDuplicate(ctx context.Context, alerts []*stockentity.StockAlert) {
	for _, alert := range alerts {
		existing, err := s.stockAlertRepo.GetAlertsByStockID(ctx, alert.StockID.String())
		duplicate := false
		if err == nil {
			for _, ea := range existing {
				if !ea.IsResolved && ea.Type == model.AlertType(alert.Type) {
					duplicate = true
					break
				}
			}
		}
		if !duplicate {
			alertModel := &model.StockAlert{}
			alertModel.FromDomain(alert)
			if err := s.stockAlertRepo.CreateAlert(ctx, alertModel); err != nil {
				fmt.Printf("Aviso: erro ao criar alerta de estoque: %v\n", err)
			}
		}
	}
}

// CreateStock cria um novo controle de estoque
func (s *Service) CreateStock(ctx context.Context, dto *stockdto.StockCreateDTO) (*stockdto.StockDTO, error) {
	// Verificar se o produto existe
	_, err := s.productRepo.GetProductById(ctx, dto.ProductID.String())
	if err != nil {
		return nil, fmt.Errorf("produto não encontrado: %w", err)
	}

	// Verificar se já existe estoque para este produto/variação
	var existingStock *model.Stock
	if dto.ProductVariationID != nil {
		existingStock, _ = s.stockRepo.GetStockByVariationID(ctx, dto.ProductVariationID.String())
	} else {
		stocks, _ := s.stockRepo.GetStockByProductID(ctx, dto.ProductID.String())
		if len(stocks) > 0 {
			existingStock = &stocks[0]
		}
	}

	if existingStock != nil {
		return nil, fmt.Errorf("já existe controle de estoque para este produto/variação")
	}

	// Criar estoque
	stock := dto.ToDomain()
	stockModel := &model.Stock{}
	stockModel.FromDomain(stock)

	if err := s.stockRepo.CreateStock(ctx, stockModel); err != nil {
		return nil, err
	}

	// Verificar alertas (com deduplicação)
	s.createAlertsIfNotDuplicate(ctx, stock.CheckAlerts())

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
	return s.stockRepo.UpdateStock(ctx, s.db, stockModel)
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

// GetStockByProductID busca estoques por produto
func (s *Service) GetStockByProductID(ctx context.Context, productID string) ([]stockdto.StockDTO, error) {
	stocksModel, err := s.stockRepo.GetStockByProductID(ctx, productID)
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

// GetStockByVariationID busca estoque por variação
func (s *Service) GetStockByVariationID(ctx context.Context, variationID string) (*stockdto.StockDTO, error) {
	stockModel, err := s.stockRepo.GetStockByVariationID(ctx, variationID)
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
	stocksModel, _, err := s.stockRepo.GetAllStocks(ctx, 0, 1000) // Fetch all (or a large page) for now as this seems to be used for reports
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
	now := time.Now().UTC()
	alert.ResolvedAt = &now

	// Fix #18: Registrar quem resolveu o alerta
	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if ok {
		userUUID := uuid.MustParse(userID)
		alert.ResolvedBy = &userUUID
	}

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

	// Calcular valor total do estoque baseado nos lotes reais
	totalStockValue := decimal.Zero
	for _, stock := range stockModels {
		batches, err := s.stockBatchRepo.GetActiveBatchesByStockID(ctx, stock.ID.String())
		if err != nil {
			continue
		}
		for _, batch := range batches {
			batchDomain := batch.ToDomain()
			totalStockValue = totalStockValue.Add(batchDomain.CurrentQuantity.Mul(batchDomain.CostPrice))
		}
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
		GeneratedAt:        time.Now().UTC(),
	}

	return report, count, nil
}

// DebitStockFromOrder debita o estoque de todos os itens de um pedido
func (s *Service) DebitStockFromOrder(ctx context.Context, orderID uuid.UUID, employeeID uuid.UUID) error {
	orderModel, err := s.orderRepo.GetOrderById(ctx, orderID.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar pedido: %w", err)
	}

	for _, groupItem := range orderModel.GroupItems {
		for _, item := range groupItem.Items {
			reason := fmt.Sprintf("Venda Pedido %s", orderID)
			quantity := decimal.NewFromFloat(item.Quantity)

			var stockModel *model.Stock

			if item.ProductVariationID != uuid.Nil {
				// Buscar estoque pela variação
				stockModel, err = s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
				if err != nil {
					continue // sem controle de estoque para esta variação, ignorar
				}
			} else {
				// Fix #23: fallback por ProductID para produtos sem variação.
				// Sem este fallback, reservas feitas pelo DebitStockFromItem pelo
				// ProductID nunca seriam debitadas via FIFO, vazando ReservedStock.
				stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
				if err != nil || len(stocks) == 0 {
					continue
				}
				stockModel = &stocks[0]
			}

			if stockModel == nil {
				continue
			}

			// Debitar FIFO
			if err := s.DebitStockFIFO(ctx, stockModel.ID, quantity, orderID, employeeID, reason); err != nil {
				return fmt.Errorf("erro ao debitar estoque para item %s: %w", item.ProductID, err)
			}

			// Verificar alertas automaticamente após cada debit de estoque (com deduplicação)
			if updatedStockModel, err := s.stockRepo.GetStockByID(ctx, stockModel.ID.String()); err == nil {
				s.createAlertsIfNotDuplicate(ctx, updatedStockModel.ToDomain().CheckAlerts())
			}
		}
	}

	return nil
}

// GetBatchesByStockID retorna todos os lotes de um estoque
func (s *Service) GetBatchesByStockID(ctx context.Context, stockID string) ([]stockdto.StockBatchDTO, error) {
	batches, err := s.stockBatchRepo.GetBatchesByStockID(ctx, stockID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar lotes do estoque: %w", err)
	}

	dtos := make([]stockdto.StockBatchDTO, 0, len(batches))
	for _, b := range batches {
		domain := b.ToDomain()
		dto := stockdto.StockBatchDTO{}
		dto.FromDomain(domain)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}
