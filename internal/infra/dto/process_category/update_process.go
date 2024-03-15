package processruledto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateProcessRuleInput struct {
	productentity.PatchProcessRule
}

func (s *UpdateProcessRuleInput) validate() error {
	if s.Name != nil && *s.Name == "" {
		return ErrNameRequired
	}

	if s.Order != nil && *s.Order < 1 {
		return ErrOrderRequired
	}

	return nil
}

func (s *UpdateProcessRuleInput) UpdateModel(model *productentity.ProcessRule) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Name != nil {
		model.Name = *s.Name
	}

	if s.Order != nil {
		model.Order = *s.Order
	}

	if s.IdealTime != nil {
		model.IdealTime = s.IdealTime
	}

	if s.ExperimentalError != nil {
		model.ExperimentalError = s.ExperimentalError
	}

	return nil
}
