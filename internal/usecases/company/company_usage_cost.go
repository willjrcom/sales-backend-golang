package companyusecases

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
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
func (s *UsageCostService) RegisterUsageCost(ctx context.Context, dto *companydto.CompanyUsageCostCreateDTO) error {
	cost, err := dto.ToDomain()
	if err != nil {
		return err
	}

	if dto.CompanyID == nil {
		companyModel, err := s.companyRepo.GetCompany(ctx)
		if err != nil {
			return err
		}
		cost.CompanyID = companyModel.ID
	}

	cost.Entity = entity.NewEntity()
	costModel := &model.CompanyUsageCost{}
	costModel.FromDomain(cost)

	return s.costRepo.Create(ctx, costModel)
}

// GetMonthlySummary returns a summary of costs for a given month
func (s *UsageCostService) GetMonthlySummary(ctx context.Context, month, year, page, perPage int) (*companydto.MonthlyCostSummaryDTO, error) {
	if s.costRepo == nil {
		return nil, ErrCostRepositoryNotSet
	}

	// Get company from context
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	// Get all costs for the month (for totals) - TODO: Optimize with aggregate query
	allCosts, err := s.costRepo.GetMonthlyCosts(ctx, companyModel.ID, month, year)
	if err != nil {
		return nil, err
	}

	// Get paginated costs for the items list
	paginatedCosts, totalItems, err := s.costRepo.GetMonthlyCostsPaginated(ctx, companyModel.ID, month, year, page, perPage)
	if err != nil {
		return nil, err
	}

	// Calculate totals by type (using allCosts)
	costsByType := make(map[string]decimal.Decimal)
	totalAmount := decimal.Zero
	totalPaid := decimal.Zero
	nfceCount := 0
	nfceCosts := decimal.Zero
	otherFee := decimal.Zero
	items := make([]companydto.CompanyUsageCostDTO, 0, len(paginatedCosts))

	for _, cost := range allCosts {
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
		} else {
			otherFee = otherFee.Add(amount)
		}

		if cost.Status == "PAID" || cost.Status == "paid" {
			totalPaid = totalPaid.Add(amount)
		}
	}

	items = make([]companydto.CompanyUsageCostDTO, 0, len(paginatedCosts))
	for _, cost := range paginatedCosts {
		item := companydto.CompanyUsageCostDTO{}
		item.FromDomain(cost.ToDomain())
		items = append(items, item)
	}

	return &companydto.MonthlyCostSummaryDTO{
		CompanyID:   companyModel.ID.String(),
		Month:       month,
		Year:        year,
		TotalAmount: totalAmount,
		TotalPaid:   totalPaid,
		CostsByType: costsByType,
		CostsCount:  len(allCosts),
		OtherFee:    otherFee,
		NFCeCosts:   nfceCosts,
		NFCeCount:   nfceCount,
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		Items:       items,
	}, nil
}
