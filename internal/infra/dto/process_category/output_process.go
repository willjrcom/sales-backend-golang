package processdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessOutput struct {
	ID uuid.UUID `json:"id"`
	productentity.ProcessCommonAttributes
}

func (s *ProcessOutput) FromModel(model *productentity.Process) {
	s.ID = model.ID
	s.Name = model.Name
	s.Order = model.Order
	s.IdealTime = model.IdealTime
	s.ExperimentalError = model.ExperimentalError
	s.CategoryID = model.CategoryID
}
