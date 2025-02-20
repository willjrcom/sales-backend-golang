package processruledto

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type ProcessRuleUpdateDTO struct {
	Name        *string `json:"name"`
	Order       *int8   `json:"order"`
	Description *string `json:"description"`
	ImagePath   *string `json:"image_path"`
	IdealTime   *string `json:"ideal_time"`
}

func (s *ProcessRuleUpdateDTO) validate() error {
	if s.Name != nil && *s.Name == "" {
		return ErrNameRequired
	}

	if s.Order != nil && *s.Order < 1 {
		return ErrOrderRequired
	}

	return nil
}

func (s *ProcessRuleUpdateDTO) UpdateDomain(processRule *productentity.ProcessRule) (err error) {
	if err = s.validate(); err != nil {
		return err
	}

	if s.Name != nil {
		processRule.Name = *s.Name
	}

	if s.Order != nil {
		processRule.Order = *s.Order
	}

	if s.Description != nil {
		processRule.Description = *s.Description
	}

	if s.ImagePath != nil {
		processRule.ImagePath = s.ImagePath
	}

	if s.IdealTime != nil {
		processRule.IdealTime, err = convertToDuration(*s.IdealTime)
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
