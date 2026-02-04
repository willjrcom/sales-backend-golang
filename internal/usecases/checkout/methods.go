package billing

import (
	"os"
	"strconv"
	"time"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

func translateMonth(m time.Month) string {
	months := []string{"", "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"}
	if m >= 1 && m <= 12 {
		return months[m]
	}
	return m.String()
}

func getEnvFloat(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func translateFrequency(p companyentity.Frequency) string {
	switch p {
	case companyentity.FrequencyMonthly:
		return "Mensal"
	case companyentity.FrequencySemiannual:
		return "Semestral"
	case companyentity.FrequencyAnnual:
		return "Anual"
	default:
		return string(p)
	}
}

func translatePlanType(p companyentity.PlanType) string {
	switch p {
	case companyentity.PlanBasic:
		return "Básico"
	case companyentity.PlanIntermediate:
		return "Intermediário"
	case companyentity.PlanAdvanced:
		return "Avançado"
	default:
		return string(p)
	}
}
