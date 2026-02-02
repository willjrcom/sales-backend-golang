package companydto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyUsageCostCreateDTO struct {
	CostType    string          `json:"cost_type"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	ReferenceID *uuid.UUID      `json:"reference_id"`
	PaymentID   *uuid.UUID      `json:"payment_id"`
	CompanyID   *uuid.UUID      `json:"company_id"`
}

func (dto *CompanyUsageCostCreateDTO) validate() error {
	if dto.CostType == "" {
		return errors.New("cost type is required")
	}
	if dto.Description == "" {
		return errors.New("description is required")
	}
	if dto.Amount.IsZero() {
		return errors.New("amount is required")
	}
	return nil
}
func (dto *CompanyUsageCostCreateDTO) ToDomain() (*companyentity.CompanyUsageCost, error) {
	if err := dto.validate(); err != nil {
		return nil, err
	}

	companyUsageCost := &companyentity.CompanyUsageCost{
		CostType:    companyentity.CostType(dto.CostType),
		Description: dto.Description,
		Amount:      dto.Amount,
		Status:      companyentity.CostStatusPending,
		ReferenceID: dto.ReferenceID,
		PaymentID:   dto.PaymentID,
	}

	if dto.CompanyID != nil {
		companyUsageCost.CompanyID = *dto.CompanyID
	}
	return companyUsageCost, nil
}
