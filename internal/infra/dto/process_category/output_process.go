package processRuledto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.ProcessRuleCommonAttributes
}

func (s *ProcessRuleOutput) FromModel(model *productentity.ProcessRule) {
	s.ID = model.ID
	s.ProcessRuleCommonAttributes = model.ProcessRuleCommonAttributes
}
