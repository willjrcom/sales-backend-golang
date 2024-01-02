package processdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateProcessInput struct {
	productentity.PatchProcess
}

func (s *UpdateProcessInput) validate() error {
	if s.Name != nil && *s.Name == "" {
		return ErrNameRequired
	}

	if s.Order != nil && *s.Order < 1 {
		return ErrOrderRequired
	}

	return nil
}

func (s *UpdateProcessInput) UpdateModel(model *productentity.Process) (err error) {
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
