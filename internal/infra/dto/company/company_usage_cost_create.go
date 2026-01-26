package companydto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyUsageCostCreateDTO struct {
	CompanyID   string          `json:"company_id"`
	CostType    string          `json:"cost_type"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
}

func (dto *CompanyUsageCostCreateDTO) ToDomain() *companyentity.CompanyUsageCost {
	return &companyentity.CompanyUsageCost{
		CostType:    companyentity.CostType(dto.CostType),
		Description: dto.Description,
		Amount:      dto.Amount,
		CompanyID:   uuid.MustParse(dto.CompanyID),
		Status:      companyentity.CostStatusPending,
	}
}
