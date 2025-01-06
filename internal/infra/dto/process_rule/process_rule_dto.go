package processruledto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleDTO struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Order             int8      `json:"order"`
	Description       string    `json:"description"`
	ImagePath         *string   `json:"image_path"`
	IdealTime         string    `json:"ideal_time"`
	ExperimentalError string    `json:"experimental_error"`
	CategoryID        uuid.UUID `json:"category_id"`
}

func (s *ProcessRuleDTO) FromDomain(model *productentity.ProcessRule) {
	*s = ProcessRuleDTO{
		ID:                model.ID,
		Name:              model.Name,
		Order:             model.Order,
		Description:       model.Description,
		ImagePath:         model.ImagePath,
		IdealTime:         getTimeFormatted(model.IdealTime),
		ExperimentalError: getTimeFormatted(model.ExperimentalError),
		CategoryID:        model.CategoryID,
	}
}

func getTimeFormatted(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}
