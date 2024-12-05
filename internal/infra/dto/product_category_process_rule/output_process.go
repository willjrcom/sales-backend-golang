package productcategoryprocessruledto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleOutput struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Order             int8      `json:"order"`
	Description       string    `json:"description"`
	ImagePath         *string   `json:"image_path"`
	IdealTime         string    `json:"ideal_time"`
	ExperimentalError string    `json:"experimental_error"`
	CategoryID        uuid.UUID `json:"category_id"`
}

func (s *ProcessRuleOutput) FromModel(model *productentity.ProcessRule) {
	s.ID = model.ID
	s.Name = model.Name
	s.Order = model.Order
	s.Description = model.Description
	s.ImagePath = model.ImagePath

	// Formatar como "HH:MM"
	s.IdealTime = getTimeFormatted(model.IdealTime)
	s.ExperimentalError = getTimeFormatted(model.ExperimentalError)
	s.CategoryID = model.CategoryID
}

func getTimeFormatted(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}
