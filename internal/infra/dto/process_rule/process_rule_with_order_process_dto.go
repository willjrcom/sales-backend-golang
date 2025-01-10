package processruledto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleWithOrderProcessDTO struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Order                 int8      `json:"order"`
	Description           string    `json:"description"`
	ImagePath             *string   `json:"image_path"`
	IdealTime             string    `json:"ideal_time"`
	CategoryID            uuid.UUID `json:"category_id"`
	TotalOrderProcessLate int       `json:"total_order_process_late"`
	TotalOrderQueue       int       `json:"total_order_queue"`
}

func (s *ProcessRuleWithOrderProcessDTO) FromDomain(processRule *productentity.ProcessRuleWithOrderProcess) {
	if processRule == nil {
		return
	}
	*s = ProcessRuleWithOrderProcessDTO{
		ID:                    processRule.ID,
		Name:                  processRule.Name,
		Order:                 processRule.Order,
		Description:           processRule.Description,
		ImagePath:             processRule.ImagePath,
		IdealTime:             getTimeFormatted(processRule.IdealTime),
		CategoryID:            processRule.CategoryID,
		TotalOrderProcessLate: processRule.TotalOrderProcessLate,
		TotalOrderQueue:       processRule.TotalOrderQueue,
	}
}
