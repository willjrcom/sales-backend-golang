package processdto

import (
	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type ProcessOutput struct {
	ID uuid.UUID `json:"id"`
	orderprocessentity.ProcessCommonAttributes
	orderprocessentity.ProcessTimeLogs
}

func (s *ProcessOutput) FromModel(model *orderprocessentity.Process) {
	s.ID = model.ID
	s.ProcessCommonAttributes = model.ProcessCommonAttributes
	s.ProcessTimeLogs = model.ProcessTimeLogs

	s.DurationFormatted = model.Duration.String()
}
