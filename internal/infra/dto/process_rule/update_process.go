package processruledto

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

	if s.Description != nil {
		model.Description = *s.Description
	}

	if s.ImagePath != nil {
		model.ImagePath = s.ImagePath
	}

	if s.IdealTime != nil {
		model.IdealTime, err = convertToDuration(*s.IdealTime)
		if err != nil {
			return err
		}
	}

	if s.ExperimentalError != nil {
		model.ExperimentalError, err = convertToDuration(*s.ExperimentalError)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertToDuration(timeStr string) (time.Duration, error) {
	// Dividir a string MM:SS
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format, expected MM:SS")
	}

	// Converter minutos e segundos para inteiros
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %v", err)
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %v", err)
	}

	// Calcular o total em nanosegundos
	duration := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

	return duration, nil
}
