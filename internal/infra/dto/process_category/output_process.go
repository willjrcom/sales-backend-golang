package processruledto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.ProcessRuleCommonAttributes
	IdealTimeFormatted         string `json:"ideal_time_formatted"`
	ExperimentalErrorFormatted string `json:"experimental_error_formatted"`
}

func (s *ProcessRuleOutput) FromModel(model *productentity.ProcessRule) {
	s.ID = model.ID
	s.ProcessRuleCommonAttributes = model.ProcessRuleCommonAttributes

	s.IdealTimeFormatted = model.IdealTime.String()
	s.ExperimentalErrorFormatted = model.ExperimentalError.String()
}
