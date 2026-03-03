package stockusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// CheckExpirations verifica todos os lotes ativos e gera alertas para os que estão vencidos ou próximos do vencimento
func (s *Service) CheckExpirations(ctx context.Context, daysThreshold int) error {
	// 1. Buscar todos os estoques ativos
	stocks, err := s.stockRepo.GetActiveStocks(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar estoques ativos: %w", err)
	}

	for _, stockModel := range stocks {
		// Fix #19: Usar GetBatchesByStockID em vez de GetActiveBatchesByStockID.
		// GetActiveBatchesByStockID filtra `expires_at > NOW()` então lotes vencidos
		// nunca chegariam ao CheckExpiration! Queremos inspecionar TODOS os lotes com estoque.
		batches, err := s.stockBatchRepo.GetBatchesByStockID(ctx, stockModel.ID.String())
		if err != nil {
			fmt.Printf("Erro ao buscar lotes para estoque %s: %v\n", stockModel.ID.String(), err)
			continue
		}

		for _, batchModel := range batches {
			batch := batchModel.ToDomain()
			alertType := batch.CheckExpiration(daysThreshold)

			if alertType != "" {
				// 3. Gerar alerta
				message := ""
				if alertType == stockentity.AlertTypeExpired {
					message = fmt.Sprintf("Lote %s do produto %s está VENCIDO (Validade: %s)",
						batchModel.ID.String()[:8], stockModel.ProductID.String()[:8], batchModel.ExpiresAt.Format("02/01/2006"))
				} else {
					message = fmt.Sprintf("Lote %s do produto %s está PRÓXIMO DO VENCIMENTO (Validade: %s)",
						batchModel.ID.String()[:8], stockModel.ProductID.String()[:8], batchModel.ExpiresAt.Format("02/01/2006"))
				}

				// Resolver ProductVariationID com segurança
				var variationID uuid.UUID
				if stockModel.ProductVariationID != nil {
					variationID = *stockModel.ProductVariationID
				}

				alert := &stockentity.StockAlert{
					Entity: entity.NewEntity(),
					StockAlertCommonAttributes: stockentity.StockAlertCommonAttributes{
						StockID:            stockModel.ID,
						Type:               alertType,
						Message:            message,
						IsResolved:         false,
						ProductID:          stockModel.ProductID,
						ProductVariationID: variationID,
					},
				}

				// Reusar helper de deduplicação do stock.go
				s.createAlertsIfNotDuplicate(ctx, []*stockentity.StockAlert{alert})
			}
		}
	}

	return nil
}

// GetExpiryAlerts retorna todos os alertas de vencimento ativos como DTOs
func (s *Service) GetExpiryAlerts(ctx context.Context) ([]stockdto.StockAlertDTO, error) {
	allAlerts, err := s.stockAlertRepo.GetActiveAlerts(ctx)
	if err != nil {
		return nil, err
	}

	var expiryAlerts []stockdto.StockAlertDTO
	for _, a := range allAlerts {
		if a.Type == model.AlertType(stockentity.AlertTypeNearExpiration) || a.Type == model.AlertType(stockentity.AlertTypeExpired) {
			alert := a.ToDomain()
			alertDTO := stockdto.StockAlertDTO{}
			alertDTO.FromDomain(alert)
			expiryAlerts = append(expiryAlerts, alertDTO)
		}
	}

	return expiryAlerts, nil
}
