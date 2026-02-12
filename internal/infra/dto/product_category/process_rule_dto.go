package productcategorydto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleDTO struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Order       int8         `json:"order"`
	Description string       `json:"description"`
	ImagePath   *string      `json:"image_path"`
	IdealTime   string       `json:"ideal_time"`
	CategoryID  uuid.UUID    `json:"category_id"`
	Category    *CategoryDTO `json:"category"`
	IsActive    bool         `json:"is_active"`
}

func (s *ProcessRuleDTO) FromDomain(processRule *productentity.ProcessRule) {
	if processRule == nil {
		return
	}

	*s = ProcessRuleDTO{
		ID:          processRule.ID,
		Name:        processRule.Name,
		Order:       processRule.Order,
		Description: processRule.Description,
		ImagePath:   processRule.ImagePath,
		IdealTime:   getTimeFormatted(processRule.IdealTime),
		CategoryID:  processRule.CategoryID,
		IsActive:    processRule.IsActive,
	}

	if processRule.Category != nil {
		category := CategoryDTO{}
		category.FromDomain(processRule.Category)
		s.Category = &category
	}
}

func getTimeFormatted(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}
