package productcategoryprocessruledto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameRequired              = errors.New("name is required")
	ErrOrderRequired             = errors.New("order min: 1is required")
	ErrIdealTimeRequired         = errors.New("ideal time is required")
	ErrExperimentalErrorRequired = errors.New("experimental error is required")
	ErrCategoryRequired          = errors.New("category ID is required")
)

type RegisterProcessRuleInput struct {
	productentity.ProcessRuleCommonAttributes
}

func (s *RegisterProcessRuleInput) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.Order < 1 {
		return ErrOrderRequired
	}

	if s.IdealTime == 0 {
		return ErrIdealTimeRequired
	}

	if s.ExperimentalError == 0 {
		return ErrExperimentalErrorRequired
	}

	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *RegisterProcessRuleInput) ToModel() (*productentity.ProcessRule, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	processRuleCommonAttributes := productentity.ProcessRuleCommonAttributes{
		Name:              s.Name,
		Order:             s.Order,
		IdealTime:         s.IdealTime * time.Second,
		ExperimentalError: s.ExperimentalError * time.Second,
		CategoryID:        s.CategoryID,
	}

	return productentity.NewProcessRule(processRuleCommonAttributes), nil
}
