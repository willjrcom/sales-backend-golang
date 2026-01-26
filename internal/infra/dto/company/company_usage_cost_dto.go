package companydto

import (
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyUsageCostDTO struct {
	ID          string          `json:"id"`
	CompanyID   string          `json:"company_id"`
	CostType    string          `json:"cost_type"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	ReferenceID *string         `json:"reference_id,omitempty"`
	PaymentID   *string         `json:"payment_id,omitempty"`
	Status      string          `json:"status"`
	CreatedAt   string          `json:"created_at"`
}

func (dto *CompanyUsageCostDTO) FromDomain(cost *companyentity.CompanyUsageCost) {
	if cost == nil {
		return
	}
	dto.ID = cost.ID.String()
	dto.CompanyID = cost.CompanyID.String()
	dto.CostType = string(cost.CostType)
	dto.Description = cost.Description
	dto.Amount = cost.Amount
	dto.CreatedAt = cost.CreatedAt.Format("2006-01-02T15:04:05Z07:00")

	if cost.ReferenceID != nil {
		refID := cost.ReferenceID.String()
		dto.ReferenceID = &refID
	}
	if cost.PaymentID != nil {
		payID := cost.PaymentID.String()
		dto.PaymentID = &payID
	}
	dto.Status = string(cost.Status)
}
