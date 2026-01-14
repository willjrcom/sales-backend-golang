package companyusecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
func (s *UsageCostService) RegisterUsageCost(ctx context.Context, companyID uuid.UUID, costType companyentity.CostType, amount decimal.Decimal, description string, referenceID *uuid.UUID) error {
	if s.costRepo == nil {
		return ErrCostRepositoryNotSet
	}

	cost := companyentity.NewCompanyUsageCost(companyID, costType, amount, description, referenceID)

	costModel := &model.CompanyUsageCost{}
	costModel.FromDomain(cost)

	return s.costRepo.Create(ctx, costModel)
}

// RegisterNFCeCost registers a cost for NFC-e emission
func (s *UsageCostService) RegisterNFCeCost(ctx context.Context, companyID uuid.UUID, invoiceID uuid.UUID, description string) error {
	cost := companyentity.NewNFCeCost(companyID, invoiceID, description)

	costModel := &model.CompanyUsageCost{}
	costModel.FromDomain(cost)

	return s.costRepo.Create(ctx, costModel)
}

// RegisterFiscalSubscriptionFee registers the monthly fiscal subscription fee (R$ 20.00)
// Only creates the cost if it doesn't already exist for the current month
func (s *UsageCostService) RegisterFiscalSubscriptionFee(ctx context.Context, companyID uuid.UUID) error {
	if s.costRepo == nil {
		return ErrCostRepositoryNotSet
	}

	// Check if fiscal subscription fee already exists for current month
	now := time.Now()
	currentMonth, currentYear := int(now.Month()), now.Year()

	costs, err := s.costRepo.GetMonthlyCosts(ctx, companyID, currentMonth, currentYear)
	if err != nil {
		return err
	}

	// Check if fiscal subscription already registered this month
	for _, cost := range costs {
		if cost.CostType == string(companyentity.CostTypeFiscalSubscription) {
			// Already registered, don't charge again
			return nil
		}
	}

	// Create new fiscal subscription cost
	cost := companyentity.NewFiscalSubscriptionCost(companyID)
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
	subscriptionFee := decimal.Zero

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
		} else if cost.CostType == string(companyentity.CostTypeSubscription) {
			subscriptionFee = subscriptionFee.Add(amount)
		}
	}

	return &companydto.MonthlyCostSummaryDTO{
		CompanyID:       companyModel.ID.String(),
		Month:           month,
		Year:            year,
		TotalAmount:     totalAmount,
		CostsByType:     costsByType,
		CostsCount:      len(costs),
		SubscriptionFee: subscriptionFee,
		NFCeCosts:       nfceCosts,
		NFCeCount:       nfceCount,
	}, nil
}

// GetCostBreakdown returns detailed breakdown of costs for a given month
func (s *UsageCostService) GetCostBreakdown(ctx context.Context, month, year int) (*companydto.CostBreakdownDTO, error) {
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

	// Convert to DTOs
	costDTOs := make([]companydto.CompanyUsageCostDTO, len(costs))
	totalAmount := decimal.Zero

	for i, cost := range costs {
		dto := companydto.CompanyUsageCostDTO{}
		dto.FromDomain(cost.ToDomain())
		costDTOs[i] = dto

		totalAmount = totalAmount.Add(cost.Amount)
	}

	return &companydto.CostBreakdownDTO{
		CompanyID:   companyModel.ID.String(),
		Month:       month,
		Year:        year,
		Costs:       costDTOs,
		TotalAmount: totalAmount.StringFixed(2),
	}, nil
}

// GetCurrentMonthSummary returns summary for the current month
func (s *UsageCostService) GetCurrentMonthSummary(ctx context.Context) (*companydto.MonthlyCostSummaryDTO, error) {
	now := time.Now()
	return s.GetMonthlySummary(ctx, int(now.Month()), now.Year())
}

// GetNextInvoicePreview returns a preview of the next billing invoice with enabled services
func (s *UsageCostService) GetNextInvoicePreview(ctx context.Context) (*companydto.NextInvoicePreviewDTO, error) {
	if s.costRepo == nil {
		return nil, ErrCostRepositoryNotSet
	}

	// Get company from context
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	company := companyModel.ToDomain()

	// Get current month costs to project usage
	now := time.Now()
	currentMonth, currentYear := int(now.Month()), now.Year()
	costs, err := s.costRepo.GetMonthlyCosts(ctx, company.ID, currentMonth, currentYear)
	if err != nil {
		return nil, err
	}

	// Calculate current month usage
	currentMonthUsage := decimal.Zero
	nfceCount := 0
	nfceCosts := decimal.Zero

	for _, cost := range costs {
		currentMonthUsage = currentMonthUsage.Add(cost.Amount)
		if cost.CostType == string(companyentity.CostTypeNFCe) {
			nfceCount++
			nfceCosts = nfceCosts.Add(cost.Amount)
		}
	}

	// Build enabled services list
	enabledServices := make([]companydto.EnabledServiceDTO, 0)

	// Fiscal subscription service (R$ 20/month when fiscal is enabled)
	fiscalSubscriptionCost := companyentity.MonthlyFiscalFee
	enabledServices = append(enabledServices, companydto.EnabledServiceDTO{
		Name:        "Assinatura Fiscal",
		Enabled:     company.FiscalEnabled,
		FixedCost:   fiscalSubscriptionCost,
		UsageCost:   decimal.Zero,
		Description: "Taxa mensal para emissão de notas fiscais",
	})

	// NFC-e service
	nfceUnitCost := companyentity.NFCeCost
	enabledServices = append(enabledServices, companydto.EnabledServiceDTO{
		Name: "NFC-e (Notas Fiscais)",

		Enabled:     company.FiscalEnabled,
		FixedCost:   decimal.Zero,
		UsageCost:   nfceCosts,
		UnitCost:    nfceUnitCost,
		UsageCount:  nfceCount,
		Description: "R$ 0,10 por nota emitida",
	})

	// Calculate next billing date (1st of next month)
	nextMonth := now.AddDate(0, 1, 0)
	nextBillingDate := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	// Calculate estimated total (only fiscal subscription if enabled, plus NFC-e costs)
	estimatedTotal := decimal.Zero
	if company.FiscalEnabled {
		estimatedTotal = fiscalSubscriptionCost.Add(nfceCosts)
	}

	return &companydto.NextInvoicePreviewDTO{
		CompanyID:         company.ID.String(),
		NextBillingDate:   nextBillingDate.Format("2006-01-02"),
		EnabledServices:   enabledServices,
		EstimatedTotal:    estimatedTotal,
		CurrentMonthUsage: currentMonthUsage,
		NFCeCount:         nfceCount,
	}, nil
}
