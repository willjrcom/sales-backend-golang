package productcategorydto

import (
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrOrderRequired     = errors.New("order min: 1is required")
	ErrIdealTimeRequired = errors.New("ideal time is required")
)

type ProcessRuleCreateDTO struct {
	Name        string    `json:"name"`
	Order       int8      `json:"order"`
	Description string    `json:"description"`
	ImagePath   *string   `json:"image_path"`
	IdealTime   string    `json:"ideal_time"`
	CategoryID  uuid.UUID `json:"category_id"`
	IsActive    *bool     `json:"is_active"`
}

func (s *ProcessRuleCreateDTO) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.Order < 1 {
		return ErrOrderRequired
	}

	if s.IdealTime == "" {
		return ErrIdealTimeRequired
	}

	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *ProcessRuleCreateDTO) ToDomain() (*productentity.ProcessRule, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	idealTime, err := convertToDuration(s.IdealTime)
	if err != nil {
		return nil, err
	}

	isActive := true
	if s.IsActive != nil {
		isActive = *s.IsActive
	}

	processRuleCommonAttributes := productentity.ProcessRuleCommonAttributes{
		Name:        s.Name,
		Order:       s.Order,
		CategoryID:  s.CategoryID,
		IdealTime:   idealTime,
		IsActive:    isActive,
		Description: s.Description,
		ImagePath:   s.ImagePath,
	}

	return productentity.NewProcessRule(processRuleCommonAttributes), nil
}
