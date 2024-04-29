package processdto

import (
	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

type ProcessOutput struct {
	ID uuid.UUID `json:"id"`
	processentity.ProcessCommonAttributes
	processentity.ProcessTimeLogs
}

func (s *ProcessOutput) FromModel(model *processentity.Process) {
	s.ID = model.ID
	s.ProcessCommonAttributes = model.ProcessCommonAttributes
	s.ProcessTimeLogs = model.ProcessTimeLogs

	s.DurationFormatted = model.Duration.String()
}
