package processdto

import (
	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type ProcessOutput struct {
	ID uuid.UUID `json:"id"`
	orderprocessentity.OrderProcessCommonAttributes
	orderprocessentity.OrderProcessTimeLogs
}

func (s *ProcessOutput) FromModel(model *orderprocessentity.OrderProcess) {
	s.ID = model.ID
	s.OrderProcessCommonAttributes = model.OrderProcessCommonAttributes
	s.OrderProcessTimeLogs = model.OrderProcessTimeLogs

	s.DurationFormatted = model.Duration.String()
}
