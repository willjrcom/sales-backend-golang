package companyusecases

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrCostRepositoryNotSet = errors.New("usage cost repository not configured")
)

type UsageCostService struct {
	costRepo    model.CompanyUsageCostRepository
	companyRepo model.CompanyRepository
}

func NewUsageCostService(costRepo model.CompanyUsageCostRepository, companyRepo model.CompanyRepository) *UsageCostService {
	return &UsageCostService{
		costRepo:    costRepo,
		companyRepo: companyRepo,
	}
}

// RegisterUsageCost registers a new usage cost for a company
func (s *UsageCostService) RegisterUsageCost(ctx context.Context, cost *companyentity.CompanyUsageCost) error {
	costModel := &model.CompanyUsageCost{}
	costModel.FromDomain(cost)

	return s.costRepo.Create(ctx, costModel)
}

// GetMonthlySummary returns a summary of costs for a given month
func (s *UsageCostService) GetMonthlySummary(ctx context.Context, month, year int) (*companydto.MonthlyCostSummaryDTO, error) {
	if s.costRepo == nil {
		return nil, ErrCostRepositoryNotSet
	}

	// Get company from context
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	// Get all costs for the month
	costs, err := s.costRepo.GetMonthlyCosts(ctx, companyModel.ID, month, year)
	if err != nil {
		return nil, err
	}

	// Calculate totals by type
	costsByType := make(map[string]decimal.Decimal)
	totalAmount := decimal.Zero
	nfceCount := 0
	nfceCosts := decimal.Zero
	otherFee := decimal.Zero

	for _, cost := range costs {
		amount := cost.Amount
		totalAmount = totalAmount.Add(amount)

		if existing, ok := costsByType[cost.CostType]; ok {
			costsByType[cost.CostType] = existing.Add(amount)
		} else {
			costsByType[cost.CostType] = amount
		}

		// Track specific metrics
		if cost.CostType == string(companyentity.CostTypeNFCe) {
			nfceCount++
			nfceCosts = nfceCosts.Add(amount)
		} else {
			otherFee = otherFee.Add(amount)
		}
	}

	return &companydto.MonthlyCostSummaryDTO{
		CompanyID:   companyModel.ID.String(),
		Month:       month,
		Year:        year,
		TotalAmount: totalAmount,
		CostsByType: costsByType,
		CostsCount:  len(costs),
		OtherFee:    otherFee,
		NFCeCosts:   nfceCosts,
		NFCeCount:   nfceCount,
	}, nil
}
